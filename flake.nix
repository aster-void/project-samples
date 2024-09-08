{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, flake-utils, ... }: flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      packages.${system}.default = pkgs.mkDerivation { };

      devShell = pkgs.mkShell {
        buildInputs = with pkgs;[
          dioxus-cli
        ];
        shellHook = ''
          hello from shellHook!
        '';
      };
    });
}
