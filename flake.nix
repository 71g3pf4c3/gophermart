{
  description = "Minimal dev shell for gophermart";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gcc
            gnumake
            git
            docker
            docker-compose
            golangci-lint
          ];

          shellHook = ''
            echo "gophermart dev shell"
            echo "available: go, gcc, make, git, docker, docker-compose, golangci-lint"
          '';
        };
      });
}
