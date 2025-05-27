{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go

    (pkgs.writeShellScriptBin "generate-examples" ''
      #!/usr/bin/env bash

      BASEPATH="${toString ./.}"
      (cd "$BASEPATH" && go install ./)
      (cd "$BASEPATH/examples" && go generate ./...)
    '')
  ];

  shellHook = ''
    export GOBIN=${toString ./bin}
    export PATH=$GOBIN:$PATH
  '';
}
