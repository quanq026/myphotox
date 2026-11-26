package server

import (
	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/api"
	"github.com/photoprism/photoprism/internal/config"
)

// APIv1 is the router group that serves the version 1 REST API.
var APIv1 *gin.RouterGroup
var registerAPIDocs func(router *gin.RouterGroup)

// registerRoutes registers the routes for handling HTTP requests with the built-in web server.
func registerRoutes(router *gin.Engine, conf *config.Config) {
	// Enables automatic redirection if the current route cannot be matched but a
	// handler for the path with (without) the trailing slash exists.
	router.RedirectTrailingSlash = true

	// Register static asset and templates routes.
	registerStaticRoutes(router, conf)

	// Register user interface routes.
	registerWebAppRoutes(router, conf)

	// Register sharing routes starting with "/s".
	registerSharingRoutes(router, conf)

	// Register ".well-known" service discovery routes.
	registerWellknownRoutes(router, conf)

	// Register JSON REST-API version 1 (APIv1) routes, grouped by functionality.
	// Docs: https://pkg.go.dev/github.com/photoprism/photoprism/internal/api

	// API Documentation.
	if registerAPIDocs != nil {
		registerAPIDocs(APIv1)
	}

	// User Sessions.
	api.CreateSession(APIv1)
	api.GetSession(APIv1)
	api.DeleteSession(APIv1)

	// Global Configuration.
	api.GetConfigOptions(APIv1)
	api.SaveConfigOptions(APIv1)
	api.StopServer(APIv1)

	// User Settings.
	api.GetClientConfig(APIv1)
	api.GetSettings(APIv1)
	api.SaveSettings(APIv1)

	// User Profile and Uploads.
	api.UploadUserFiles(APIv1)
	api.ProcessUserUpload(APIv1)
	api.UploadUserAvatar(APIv1)
	api.FindUserSessions(APIv1)
	api.CreateUserPasscode(APIv1)
	api.ConfirmUserPasscode(APIv1)
	api.ActivateUserPasscode(APIv1)
	api.DeactivateUserPasscode(APIv1)
	api.UpdateUserPassword(APIv1)
	api.UpdateUser(APIv1)

	// Thumbnail Images.
	api.GetThumb(APIv1)

	// Video Streaming.
	api.GetVideo(APIv1)

	// Downloads.
	api.GetDownload(APIv1)
	api.ZipCreate(APIv1)
	api.ZipDownload(APIv1)

	// Index and Import.
	api.StartImport(APIv1)
	api.CancelImport(APIv1)
	api.StartIndexing(APIv1)
	api.CancelIndexing(APIv1)

	// Photo Search and Organization.
	api.SearchPhotos(APIv1)
	api.SearchPhotosView(APIv1)
	api.GetPhoto(APIv1)
	api.UpdatePhoto(APIv1)
	api.GetPhotoDownload(APIv1)
	api.ApprovePhoto(APIv1)
	api.LikePhoto(APIv1)
	api.DislikePhoto(APIv1)
	api.AddPhotoLabel(APIv1)
	api.RemovePhotoLabel(APIv1)
	api.UpdatePhotoLabel(APIv1)
	api.GetFile(APIv1)
	api.GetFileBytes(APIv1)
	api.DeleteFile(APIv1)
	api.ChangeFileOrientation(APIv1)
	api.PhotoPrimary(APIv1)
	api.PhotoUnstack(APIv1)

	// Photo Albums (includes AlbumLinks for public sharing).
	api.SearchAlbums(APIv1)
	api.GetAlbum(APIv1)
	api.AlbumCover(APIv1)
	api.CreateAlbum(APIv1)
	api.UpdateAlbum(APIv1)
	api.DeleteAlbum(APIv1)
	api.DownloadAlbum(APIv1)
	api.GetAlbumLinks(APIv1)
	api.CreateAlbumLink(APIv1)
	api.UpdateAlbumLink(APIv1)
	api.DeleteAlbumLink(APIv1)
	api.LikeAlbum(APIv1)
	api.DislikeAlbum(APIv1)
	api.CloneAlbums(APIv1)
	api.AddPhotosToAlbum(APIv1)
	api.RemovePhotosFromAlbum(APIv1)

	// Files and Folders.
	api.SearchFoldersOriginals(APIv1)
	api.SearchFoldersImport(APIv1)
	api.FolderCover(APIv1)

	// Batch Operations.
	api.BatchAlbumsDelete(APIv1)
	api.BatchPhotosEdit(APIv1)
	api.BatchPhotosApprove(APIv1)
	api.BatchPhotosArchive(APIv1)
	api.BatchPhotosRestore(APIv1)
	api.BatchPhotosPrivate(APIv1)
	api.BatchPhotosDelete(APIv1)

	// Technical Endpoints.
	api.GetSvg(APIv1)
	api.GetStatus(APIv1)
	api.GetErrors(APIv1)
	api.DeleteErrors(APIv1)
	api.WebSocket(APIv1)
	api.Options(APIv1)
	api.Echo(APIv1)
}
