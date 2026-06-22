package report

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/openclaw/discrawl/internal/store"
)

func TestReportsCompareLegacyVariableWidthTimestampsChronologically(t *testing.T) {
	ctx := context.Background()
	s, err := store.Open(ctx, filepath.Join(t.TempDir(), "discrawl.db"))
	require.NoError(t, err)
	defer func() { _ = s.Close() }()

	boundary := time.Date(2026, 4, 25, 12, 0, 0, 0, time.UTC)
	halfSecond := boundary.Add(500 * time.Millisecond)
	require.NoError(t, s.UpsertGuild(ctx, store.GuildRecord{ID: "g1", Name: "Guild", RawJSON: `{}`}))
	require.NoError(t, s.UpsertChannel(ctx, store.ChannelRecord{ID: "c1", GuildID: "g1", Kind: "text", Name: "boundary", RawJSON: `{}`}))
	require.NoError(t, s.UpsertChannel(ctx, store.ChannelRecord{ID: "c2", GuildID: "g1", Kind: "text", Name: "half-only", RawJSON: `{}`}))
	require.NoError(t, s.UpsertChannel(ctx, store.ChannelRecord{ID: "c3", GuildID: "g1", Kind: "text", Name: "never", RawJSON: `{}`}))

	require.NoError(t, s.UpsertMessages(ctx, []store.MessageMutation{
		{
			Record: store.MessageRecord{ID: "m1", GuildID: "g1", ChannelID: "c1", ChannelName: "boundary", AuthorID: "u1", AuthorName: "u1", CreatedAt: boundary.Format(time.RFC3339Nano), Content: "exact", NormalizedContent: "exact", RawJSON: `{}`},
		},
		{
			Record:   store.MessageRecord{ID: "m2", GuildID: "g1", ChannelID: "c1", ChannelName: "boundary", AuthorID: "u2", AuthorName: "u2", CreatedAt: halfSecond.Format(time.RFC3339Nano), Content: "half", NormalizedContent: "half", RawJSON: `{}`},
			Mentions: []store.MentionEventRecord{{MessageID: "m2", GuildID: "g1", ChannelID: "c1", AuthorID: "u2", TargetType: "user", TargetID: "u1", TargetName: "u1", EventAt: halfSecond.Format(time.RFC3339Nano)}},
		},
		{
			Record: store.MessageRecord{ID: "m3", GuildID: "g1", ChannelID: "c2", ChannelName: "half-only", AuthorID: "u3", AuthorName: "u3", CreatedAt: halfSecond.Format(time.RFC3339Nano), Content: "half only", NormalizedContent: "half only", RawJSON: `{}`},
		},
	}))

	_, err = s.DB().ExecContext(ctx, `
		update messages
		set created_at = case id
			when 'm1' then ?
			when 'm2' then ?
			when 'm3' then ?
		end
		where id in ('m1', 'm2', 'm3')
	`, boundary.Format(time.RFC3339Nano), halfSecond.Format(time.RFC3339Nano), halfSecond.Format(time.RFC3339Nano))
	require.NoError(t, err)
	_, err = s.DB().ExecContext(ctx, `update mention_events set event_at = ? where message_id = 'm2'`, halfSecond.Format(time.RFC3339Nano))
	require.NoError(t, err)

	activity, err := Build(ctx, s, Options{Now: boundary.Add(time.Second)})
	require.NoError(t, err)
	require.Equal(t, halfSecond, activity.LatestMessageAt)

	digest, err := BuildDigest(ctx, s, DigestOptions{Now: boundary.Add(time.Second), Since: time.Second})
	require.NoError(t, err)
	require.Equal(t, 3, digest.Totals.Messages)
	require.Equal(t, "c1", digest.Channels[0].ChannelID)
	require.Equal(t, 2, digest.Channels[0].Messages)
	require.Equal(t, 1, digest.Channels[0].TopMentions[0].Count)

	quiet, err := BuildQuiet(ctx, s, QuietOptions{Now: boundary.Add(2 * time.Second), Since: 2 * time.Second})
	require.NoError(t, err)
	require.Len(t, quiet.Channels, 1)
	require.Equal(t, "c3", quiet.Channels[0].ChannelID)

	trends, err := BuildTrends(ctx, s, TrendsOptions{Now: boundary, Weeks: 1, GuildID: "g1"})
	require.NoError(t, err)
	require.Len(t, trends.Rows, 3)
	for _, row := range trends.Rows {
		require.Equal(t, 0, row.Weekly[0].Messages)
	}
}
