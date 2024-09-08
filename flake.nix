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
      packages.x86_64-linux.hello = pkgs.hello;
      packages.x86_64-linux.default = pkgs.hello;
      devShell = pkgs.mkShell {
        buildInputs = [ pkgs.hello ];
        shellHook = ''
          hello from shellHook!
        '';
      };
    });
}
