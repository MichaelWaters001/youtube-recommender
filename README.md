# YouTube Recommender

YouTube Recommender is a Go-based web service that helps users discover YouTube content creators based on user-submitted tags and community votes. Users can log in via Google OAuth, submit tags for their favorite creators, and vote on existing tags to refine recommendations. The service provides fast search functionality to find creators by exact tag matches.

## Features

- Google OAuth Authentication – Secure login using Google accounts.
- Tagging System – Users can tag YouTube creators with relevant topics.
- Voting Mechanism – Community-driven ranking via upvotes/downvotes.
- YouTube API Integration – Automatically fetches creator details by YouTube ID.
- JWT-Based Authentication – Stateless and secure token-based login.
- Fast Search API – Find creators by exact tag matches.

## Tech Stack

- Backend: Go (Golang)
- Web Framework: Gin
- Database: PostgreSQL
- Authentication: Google OAuth & JWT
- Logging: slog

## Setup Instructions

### Install Dependencies
Make sure you have Go and PostgreSQL installed.

```sh
git clone https://github.com/MichaelWaters001/youtube-recommender.git
cd youtube-recommender
go mod tidy
```

---

### Configure the Database
Set up your PostgreSQL database:

```sh
psql -U your_db_user -d youtube_recommender -f migrations/001_init.sql
```

Edit `config.toml` to match your database credentials:

```toml
[database]
host = "localhost"
port = 5432
user = "your_db_user"
password = "your_db_password"
dbname = "youtube_recommender"
sslmode = "disable"
```

---

### Set Up OAuth Credentials
Get your Google OAuth Client ID & Secret from [Google Developers Console](https://console.developers.google.com/) and add them to `config.toml`:

```toml
[oauth]
client_id = "your-google-client-id"
client_secret = "your-google-client-secret"
redirect_url = "http://localhost:8080/auth/google/callback"
```

---

### Run the Server
```sh
go run cmd/main.go
```

The API will be available at `http://localhost:8080`.

---

## API Endpoints & Examples

### Authentication
| Method | Endpoint                   | Description |
|--------|----------------------------|-------------|
| `GET`  | `/auth/google`             | Redirects to Google OAuth login |
| `GET`  | `/auth/google/callback`    | Handles OAuth callback and returns JWT |
| `POST` | `/auth/logout`             | Logs out the user |

#### Example: Google OAuth Login
```sh
curl -X GET http://localhost:8080/auth/google
```

#### Example: Logout
```sh
curl -X POST http://localhost:8080/auth/logout
```

---

### Creators
| Method  | Endpoint                | Description |
|---------|-------------------------|-------------|
| `POST`  | `/creators`             | Add a new YouTube creator |
| `GET`   | `/creators/:id`         | Get creator details |

#### Example: Add a Creator
```sh
curl -X POST http://localhost:8080/creators \
     -H "Content-Type: application/json" \
     -d '{"youtube_id": "UC123456"}'
```

#### Example: Get Creator by ID
```sh
curl -X GET http://localhost:8080/creators/1
```

---

### Tags
| Method  | Endpoint                        | Description |
|---------|---------------------------------|-------------|
| `POST`  | `/creators/:id/tags`           | Add a tag to a creator |
| `GET`   | `/creators/:id/tags`           | Get tags for a creator |

#### Example: Add a Tag to a Creator
```sh
curl -X POST http://localhost:8080/creators/1/tags \
     -H "Content-Type: application/json" \
     -d '{"tag_name": "Tech"}' \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Example: Get Tags for a Creator
```sh
curl -X GET http://localhost:8080/creators/1/tags
```

---

### Voting
| Method  | Endpoint                   | Description |
|---------|----------------------------|-------------|
| `POST`  | `/votes`                   | Upvote or downvote a tag |
| `DELETE`| `/votes/:creator_tag_id`    | Remove a vote |

#### Example: Upvote a Tag
```sh
curl -X POST http://localhost:8080/votes \
     -H "Content-Type: application/json" \
     -d '{"creator_tag_id": 1, "vote_type": 1}' \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Example: Downvote a Tag
```sh
curl -X POST http://localhost:8080/votes \
     -H "Content-Type: application/json" \
     -d '{"creator_tag_id": 1, "vote_type": -1}' \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Example: Remove a Vote
```sh
curl -X DELETE http://localhost:8080/votes/1 \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### Search
| Method  | Endpoint                   | Description |
|---------|----------------------------|-------------|
| `GET`   | `/search?tag=example`      | Search for creators by tag |

#### Example: Search for Creators by Tag
```sh
curl -X GET "http://localhost:8080/search?tag=Tech"
```

---

## License
This project is open-source and available under the MIT License.