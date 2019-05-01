class CliV2 < Formula
  desc "Cli v2 (using Go)"
  homepage "https://github.com/Qordobacode/Cli-v2"
  url "https://github.com/Qordobacode/Cli-v2/archive/version-0.1.tar.gz"
  sha256 "aafa200e49341bfbe0b0ebc7c7ab53df606da976243f37b688dc0136d1215346"

  depends_on "go" => :build

  def install
    ENV["GOPATH"] = buildpath
    system "build.sh"
    bin.install ".gobuild/bin/cli-v2" => "qor"
  end

  test do
    system "#{bin}/qor", "--version"
  end
end