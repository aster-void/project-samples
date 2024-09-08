{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        name = "go-sample-flake";
      in
      rec {
        packages.${name} = pkgs.buildGo122Module {
          inherit name;
          pname = name;
          src = ./.;
          vendorHash = null;
        };
        packages.default = packages.${name};

        devShell = pkgs.mkShell {
          inherit name;
          buildInputs = [ ];
          shellHook = '''';
        };
      });
}
