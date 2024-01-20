package config

func GetOrigins() string {
	origins := []string{
		"http://localhost:3000",
		"https://git-leet-frontend.vercel.app",
		"http://localhost:5173",
	}

	allowedOriginsString := ""

	for i, origin := range origins {
		if i == len(origins)-1 {
			allowedOriginsString += origin
		} else {
			allowedOriginsString += (origin + ", ")
		}
	}

	return allowedOriginsString
}
