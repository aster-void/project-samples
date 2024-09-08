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
        name = "python-sample";
        pyvenv = pkgs.python310.withPackages (ps: [
          ps.pip
          ps.virtualenv
        ]);
      in
      {
        packages.default = pkgs.stdenv.mkDerivation {
          inherit name;
          src = ./.;
          nativeBuildInputs = with pkgs; [ python312 ];
          # configure buildPhase as you need
          buildPhase = ''
            mkdir -p $out/bin
            cat > $out/bin/${name} << EOF
            #!/usr/bin/env python
            print("Hello flake-utils via ${system}!")
            EOF
            chmod u+x $out/bin/${name}
          '';
        };

        devShell = pkgs.mkShell
          {
            inherit name;
            buildInputs = [ ];
          } // {
          # this is necessary when python acts as a wrapper of C/C++ (e.g. when you use TF, NP, etc)
          LD_LIBRARY_PATH = "${pkgs.stdenv.cc.cc.lib}/lib";

          # these are required if you use venv+requirements.txt
          packages = [ pyvenv ];
          shellHook = ''
            python -m venv ./.venv
            source ./.venv/bin/activate
            pip install -r requirements.txt
          '';
        };
      });
}
