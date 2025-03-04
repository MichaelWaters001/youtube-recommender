import { useState } from "react";
import { Chip, IconButton, Typography, Grid } from "@mui/material";
import ThumbUpIcon from "@mui/icons-material/ThumbUp";
import ThumbDownIcon from "@mui/icons-material/ThumbDown";
import { voteTag, removeVote } from "../api/tags";

function TagList({ tags }) {
  const [tagVotes, setTagVotes] = useState(tags);

  const handleVote = async (creatorTagId, voteType) => {
    const res = await voteTag(creatorTagId, voteType);
    if (res.error) {
      alert(res.error);
      return;
    }
    setTagVotes(
      tagVotes.map((tag) =>
        tag.id === creatorTagId
          ? { ...tag, upvotes: voteType === 1 ? tag.upvotes + 1 : tag.upvotes }
          : tag
      )
    );
  };

  return (
    <Grid container spacing={1}>
      {tagVotes.map((tag) => (
        <Grid item key={tag.id}>
          <Chip label={`${tag.name} (${tag.upvotes}👍)`} variant="outlined" />
          <IconButton onClick={() => handleVote(tag.id, 1)}>
            <ThumbUpIcon />
          </IconButton>
          <IconButton onClick={() => handleVote(tag.id, -1)}>
            <ThumbDownIcon />
          </IconButton>
        </Grid>
      ))}
    </Grid>
  );
}

export default TagList;
