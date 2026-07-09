# Repo guidelines
## Release process
The release process is automated using **goreleaser** which works in *snapshot* mode on the `develop` branch in order to test the release process.    
Meanwhile, the release process is triggered by assigning a **tag** in a commit where the tag name is in the format `vX.Y.Z`, following the semantic versioning. So you can create a release in any branch, but it is recommended to create a release from the `main` branch or to merge into `main` right after the release.