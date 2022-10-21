package tmdb

const BASEURL = "https://api.themoviedb.org/3"
const MoviesTopRatedPath = "movie/top_rated"

var timeoutMap = map[string]int{
	MoviesTopRatedPath: 1000,
}
