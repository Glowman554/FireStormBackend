package remote

import (
	"context"
	"errors"

	"encore.dev/beta/auth"
)

func isValidName(name string) bool {
	for _, c := range name {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_') {
			return false
		}
	}
	return true
}

//encore:api auth method=GET path=/package/create/:pkgName
func CreatePackage(ctx context.Context, pkgName string) error {
	if !isValidName(pkgName) {
		return errors.New("invalid package name")
	}

	uid, ok := auth.UserID()
	if !ok {
		return errors.New("missing uid")
	}
	pkg := Package{
		Owner:   string(uid),
		Package: pkgName,
	}

	err := createPackage(ctx, pkg)
	if err != nil {
		return err
	}

	_, err = PackageCreation.Publish(ctx, &pkg)
	return err
}

//encore:api auth method=GET path=/package/delete/full/:pkgName
func DeletePackage(ctx context.Context, pkgName string) error {
	uid, ok := auth.UserID()
	if !ok {
		return errors.New("missing uid")
	}

	pkg, err := loadPackage(ctx, pkgName)
	if err != nil {
		return err
	}

	if pkg.Owner != string(uid) {
		return errors.New("you are not allowed to delete this package")
	}
	err = deletePackage(ctx, pkgName)
	if err != nil {
		return err
	}

	_, err = PackageDeletion.Publish(ctx, pkg)
	return err
}

//encore:api auth method=GET path=/package/delete/version/:pkgName/:versionName
func DeletePackageVersion(ctx context.Context, pkgName string, versionName string) error {
	uid, ok := auth.UserID()
	if !ok {
		return errors.New("missing uid")
	}

	pkg, err := loadPackage(ctx, pkgName)
	if err != nil {
		return err
	}

	if pkg.Owner != string(uid) {
		return errors.New("you are not allowed to delete this package version")
	}
	return deletePackageVersion(ctx, pkgName, versionName)
}

type UploadFileProps struct {
	Content string `json:"content"`
}

//encore:api auth method=POST path=/file/upload/:pkgName/:versionName/*file
func UploadFile(ctx context.Context, pkgName string, versionName string, file string, props *UploadFileProps) error {
	uid, ok := auth.UserID()
	if !ok {
		return errors.New("missing uid")
	}

	pkg, err := loadPackage(ctx, pkgName)
	if err != nil {
		return err
	}

	if pkg.Owner != string(uid) {
		return errors.New("you are not allowed to upload to this package")
	}

	updateDateUpdated(ctx, pkg.Package)

	return createFile(ctx, File{
		Package: *pkg,
		Name:    file,
		Version: versionName,
		Content: props.Content,
	})
}

type ListFilesResponse struct {
	Files []string `json:"files"`
}

//encore:api public method=GET path=/package/list/:pkgName/:versionName
func ListFiles(ctx context.Context, pkgName string, versionName string) (*ListFilesResponse, error) {
	files, err := loadVersionFileList(ctx, versionName, pkgName)
	if err != nil {
		return nil, err
	}
	return &ListFilesResponse{
		Files: files,
	}, nil
}

type LoadFileResponse struct {
	Content string `json:"content"`
}

//encore:api public method=GET path=/package/load/file/:pkgName/:versionName/*file
func LoadFile(ctx context.Context, pkgName string, versionName string, file string) (*LoadFileResponse, error) {
	content, err := loadFile(ctx, file, versionName, pkgName)
	if err != nil {
		return nil, err
	}
	return &LoadFileResponse{
		Content: *content,
	}, nil
}

type ListPackagesProps struct {
	Limit  int `json:"limit" query:"limit"`
	Offset int `json:"offset" query:"offset"`
}
type ListPackagesResponse struct {
	Packages []Package `json:"packages"`
}

//encore:api public method=GET path=/package/all
func ListPackages(ctx context.Context, props *ListPackagesProps) (*ListPackagesResponse, error) {
	if props.Limit > 10 {
		props.Limit = 10
	}
	pkgs, err := loadPackages(ctx, props.Limit, props.Offset)
	if err != nil {
		return nil, err
	}
	return &ListPackagesResponse{
		Packages: pkgs,
	}, nil
}

//encore:api public method=GET path=/package/get/:pkgName
func GetPackage(ctx context.Context, pkgName string) (*Package, error) {
	pkg, err := loadPackage(ctx, pkgName)
	if err != nil {
		return nil, err
	}
	return pkg, nil

}

type GetVersionsResponse struct {
	Versions []string `json:"versions" query:"version"`
}

//encore:api public method=GET path=/package/version/list/:pkgName
func GetVersions(ctx context.Context, pkgName string) (*GetVersionsResponse, error) {
	versions, err := loadPackageVersions(ctx, pkgName)
	if err != nil {
		return nil, err
	}

	return &GetVersionsResponse{
		Versions: versions,
	}, nil
}
