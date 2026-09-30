package clieant

import (
	"context"
	"encoding/json"
	"fmt"
	"homework-2/internal/movie"
	"net/http"
	"time"
)

func Get(id int, timeout time.Duration, ctx context.Context) (movie.Movie, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://homeworksite.site/%d/info.0.json", id),
		nil,
	)
	if err != nil {
		return movie.Movie{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return movie.Movie{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return movie.Movie{}, fmt.Errorf("movie %d: status %d", id, res.StatusCode)
	}

	var data movie.Movie
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return movie.Movie{}, fmt.Errorf("movie %d: invalid json", id)
	}

	return data, nil
}
