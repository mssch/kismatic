# Vendoring Ansible

Using python's `pip`, it is possible to install a python package and all it's dependencies into a specific directory. 
Taking advantage of this capability, we can vendor Ansible and it's dependencies. The only caveat is that the `PYTHONPATH`
has to point to the vendored folder. 

## Using Docker to create vendored package
The Dockerfile creates an image that contains all the dependencies that are necessary to install Ansible using pip.

Once the image is built, run the following to obtain a vendored Ansible package
```
docker run --rm -v $(pwd)/out:/ansible apprenda/vendor-ansible \ 
    pip install --install-option="--prefix=/ansible" ansible
```

The `--install-option` flag is telling pip to install the ansible package to the `/ansible` directory.
The `/ansible` directory is in turn a Docker volume that maps to `$(pwd)/out`, which will contain the vendored bits:

## Use vendored ansible
In order to use Ansible, the `PYTHONPATH` must include the `lib` and `lib64` directories. Make sure to use absolute paths.
```
PYTHONPATH=$(pwd)/out/lib/python2.7/site-packages/:$(pwd)/out/lib64/python2.7/site-packages bin/ansible
```