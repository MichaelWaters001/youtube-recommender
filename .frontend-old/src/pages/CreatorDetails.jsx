import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Container, Typography, Card, CardContent, CardMedia } from "@mui/material";

function CreatorDetails() {
  const { id } = useParams(); // Get creator ID from URL
  const [creator, setCreator] = useState(null);

  useEffect(() => {
    fetch(`/api/creators/${id}`)
      .then((res) => res.json())
      .then((data) => setCreator(data))
      .catch((err) => console.error("Failed to fetch creator", err));
}, [id]);

  if (!creator) return <p>Loading...</p>;

  return (
    <Container>
      <Card sx={{ maxWidth: 800, margin: "auto", mt: 5 }}>
        <CardMedia
          component="img"
          height="200"
          image={`https://img.youtube.com/vi/${creator.youtube_id}/hqdefault.jpg`} // Display thumbnail
          alt={creator.name}
        />
        <CardContent>
          <Typography variant="h4" gutterBottom>{creator.name}</Typography>
          <Typography variant="body1" color="textSecondary">{creator.description}</Typography>
        </CardContent>
      </Card>
    </Container>
  );
}

export default CreatorDetails;