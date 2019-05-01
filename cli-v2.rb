# Documentation: https://docs.brew.sh/Formula-Cookbook
#                https://rubydoc.brew.sh/Formula
# PLEASE REMOVE ALL GENERATED COMMENTS BEFORE SUBMITTING YOUR PULL REQUEST!
class CliV2 < Formula
  desc "Cli v2 (using Go)"
  homepage "https://github.com/Qordobacode/Cli-v2"
  url "https://github.com/Qordobacode/Cli-v2/releases/tag/version-0.2"
  sha256 "3246fb72b44e5940dde54c28dc8484d7bd1c7be9d03e0d273d338539ba469dce"
  depends_on "go" => :build

  def install
     system "gobuild.sh"
     bin.install ".gobuild/bin/cli-v2" => "qor"
  end

  test do
    system "#{bin}/qor", "--help"
  end
end
