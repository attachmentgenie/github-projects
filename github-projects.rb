# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class GithubProjects < Formula
  desc ""
  homepage ""
  version "0.5.0"
  depends_on :macos

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/attachmentgenie/github-projects/releases/download/v0.5.0/github-projects_0.5.0_darwin_arm64.tar.gz"
      sha256 "7bbd33ac939036588b6aa9f1301ad098dc34a94d03614a4590e94577021e4ef3"

      def install
        bin.install "github-projects"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/attachmentgenie/github-projects/releases/download/v0.5.0/github-projects_0.5.0_darwin_amd64.tar.gz"
      sha256 "2b76b26f4711010b6e3bdbfed829e7ac77358e2a19e5c848085238e679acd18a"

      def install
        bin.install "github-projects"
      end
    end
  end
end
