let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-23.11";
  pkgs = import nixpkgs { config = {}; overlays = []; };
in

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    git
    gopls
    lunarvim
    neovim
    zellij
  ];

  shellHook = ''
    echo "Welcome to your nix-shell environment!"
    
    # Setting aliases
    alias vim=lvim
    
    # Launch zellij
    zellij --layout layout.kdl options --simplified-ui true --default-shell zsh --copy-command "xclip -selection clipboard"
  '';
}
