{ pkgs ? import <nixpkgs> {} }:

with pkgs;

buildGoModule rec {
  pname = "ska";
  version = "latest";

  src = ./.;

  vendorSha256 = "3nk161ayQDfS39TuDXfnuOmept8ya2kT9tkP/Cij7Jc=";

  nativeBuildInputs = [ installShellFiles ];

  postInstall = ''
    installShellCompletion --zsh completions/ska.zsh
    installShellCompletion --fish completions/ska.fish

    mkdir $out/share/man/
    cp $src/man/ska.1 $out/share/man/
  '';
}
