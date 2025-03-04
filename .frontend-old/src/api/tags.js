const API_URL = "http://localhost:8080";

export const voteTag = async (creatorTagId, voteType) => {
  const token = localStorage.getItem("token");
  if (!token) return { error: "User not logged in" };

  const res = await fetch(`${API_URL}/votes`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ creator_tag_id: creatorTagId, vote_type: voteType }),
  });

  return res.json();
};

export const removeVote = async (creatorTagId) => {
  const token = localStorage.getItem("token");
  if (!token) return { error: "User not logged in" };

  const res = await fetch(`${API_URL}/votes/${creatorTagId}`, {
    method: "DELETE",
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return res.json();
};
