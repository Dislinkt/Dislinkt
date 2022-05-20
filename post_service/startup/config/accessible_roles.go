package config

// => /naziv_paketa.naziv_servisa/naziv_metode : { naziv_role1, naziv_role2 }
// za sve metode koje ne treba da se presrecu -> ne dodaju se u mapu
func AccessibleRoles() map[string][]string {
	const postService = "/post_service_proto.PostService/"

	return map[string][]string{
		postService + "CreatePost":    {"Regular"},
		postService + "CreateComment": {"Regular"},
		postService + "LikePost":      {"Regular"},
		postService + "DislikePost":   {"Regular"},
	}
}
