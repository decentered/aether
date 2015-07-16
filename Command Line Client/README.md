## Command Line Client

It's a pure Python package with some system dependencies. It does not depend on rest of the repository, you can grab the command line folder on its own. Should run on any UNIX (Linux, OS X, whatever else). I have no idea about Windows.

Do ```pip install -r requirements.txt``` and follow through with the dependencies. Then, ```python main.py```. If you want to keep it running after you disconnect the ssh connection, use tmux.

Keep in mind that the next major version of Aether will be changing the protocol. So this CLI release does not guarantee API stability. Caveat emptor!

License: same as the repo. (AGPL)