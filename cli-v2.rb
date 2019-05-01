class CliV2 < Formula
  homepage "https://github.com/Qordobacode/Cli-v2"
  url "https://github.com/Qordobacode/Cli-v2/archive/version-0.1.tar.gz"
  sha256 "8e4bd52a7204526c7ce649bf954c576e0e15c24adcb52e0a1edcdb1204996c79"

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