import { useState, useEffect } from "react";
import { Container, Grid, Typography, TextField, Card, CardContent, Button } from "@mui/material";
import { Link } from "react-router-dom";

function Home() {
  const [creators, setCreators] = useState([]);
  const [search, setSearch] = useState("");

  useEffect(() => {
    fetch("http://localhost:8080/creators")
      .then((res) => res.json())
      .then((data) => setCreators(data))
      .catch((err) => console.error("Failed to fetch creators", err));
  }, []);

  return (
    <Container>
      <Typography variant="h4" gutterBottom>YouTube Creators</Typography>
      <TextField
        fullWidth
        label="Search by tag"
        variant="outlined"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        sx={{ mb: 3 }}
      />
      <Grid container spacing={3}>
        {creators.map((creator) => (
          <Grid item xs={12} sm={6} md={4} key={creator.id}>
            <Card>
              <CardContent>
                <Typography variant="h6">{creator.name}</Typography>
                <Button component={Link} to={`/creators/${creator.id}`} variant="contained" color="primary" sx={{ mt: 2 }}>
                  View Details
                </Button>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Container>
  );
}

export default Home;
