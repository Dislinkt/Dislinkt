package config

func AccessiblePermissions() map[string]string {
	const postService = "/post_service_proto.PostService/"

	return map[string]string{
		postService + "createPost":     "createPostPermission",
		postService + "CreateComment":  "createCommentPermission",
		postService + "LikePost":       "likePostPermission",
		postService + "DislikePost":    "dislikePostPermission",
		postService + "getRecent":      "getRecentPermission",
		postService + "CreateJobOffer": "createJobOffer",
	}
}
