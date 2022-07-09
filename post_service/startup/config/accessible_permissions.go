package config

func AccessiblePermissions() map[string]string {
	const postService = "/post_service_proto.PostService/"

	return map[string]string{
		postService + "createPost":     "createPostPermission",
		postService + "createComment":  "createCommentPermission",
		postService + "likePost":       "likePostPermission",
		postService + "dislikePost":    "dislikePostPermission",
		postService + "getRecent":      "getRecentPermission",
		postService + "CreateJobOffer": "createJobOffer",
	}
}
