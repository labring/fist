set nocompatible              " be iMproved, required
set backspace=2
let mapleader=","
filetype off "required
set foldmethod=indent
set encoding=utf-8
set rtp+=~/.vim/bundle/Vundle.vim
call vundle#begin()
Plugin 'VundleVim/Vundle.vim'
Plugin 'L9'
Plugin 'fatih/vim-go'
Bundle 'scrooloose/nerdtree'
Plugin 'Valloric/YouCompleteMe'
Plugin 'vim-scripts/taglist.vim'
Plugin 'tomasr/molokai'
Plugin 'SirVer/ultisnips'
Plugin 'honza/vim-snippets'
Plugin 'andrewstuart/vim-kubernetes'
Plugin 'majutsushi/tagbar'
call vundle#end()            " required
nmap <C-h> <C-w>h
nmap <C-l> <C-w>l
nmap <C-j> <C-w>j
nmap <C-k> <C-w>k
map <C-n> :NERDTreeToggle<CR>
map <C-t> :b 1 <CR>
nmap <Leader>t :TagbarToggle<CR>
let g:tagbar_width = 25
autocmd VimEnter * nested :call tagbar#autoopen(1)
autocmd BufWritePre *.go:Fmt
filetype plugin indent on    " required

let g:tagbar_right = 1
let g:molokai_original = 1
set sw=4

set nu
set ts=4
"set expandtab
syntax on

let g:go_highlight_functions = 1
let g:go_highlight_methods = 1
let g:go_highlight_structs = 1
let g:go_highlight_operators = 1
let g:go_highlight_build_constraints = 1
let g:go_fmt_command = "goimports"
let g:go_fmt_fail_silently = 1
let g:syntastic_go_checkers = ['golint', 'govet', 'errcheck']
let g:syntastic_mode_map = { 'mode': 'active', 'passive_filetypes': ['go'] }

let g:UltiSnipsExpandTrigger="<C-g>"
let g:godef_split=2
let g:godef_same_file_in_same_window=1

au FileType go nmap <leader>r <Plug>(go-run)
au FileType go nmap <leader>b <Plug>(go-build)
au FileType go nmap <leader>T <Plug>(go-test)
au FileType go nmap <leader>c <Plug>(go-coverage)

au FileType go nmap <Leader>ds <Plug>(go-def-split)
au FileType go nmap <Leader>dv <Plug>(go-def-vertical)
au FileType go nmap <Leader>dt <Plug>(go-def-tab)

let g:syntastic_always_populate_loc_list = 1
let g:syntastic_auto_loc_list = 1
let g:syntastic_check_on_open = 1
let g:syntastic_check_on_wq = 0

set rtp+=$GOPATH/src/github.com/golang/lint/misc/vim
autocmd BufWritePost,FileWritePost *.go execute 'Lint' | cwindow
set completefunc=emoji#complete
map <C-e> :%s/:\([^:]\+\):/\=emoji#for(submatch(1), submatch(0))/g <CR>

set lines=120 columns=150
