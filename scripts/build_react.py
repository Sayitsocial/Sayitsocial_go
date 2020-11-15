import os
from os import getcwd, system
import subprocess
import shutil

react_path = os.path.abspath(os.path.join(
    os.path.split(__file__)[0], "../web/v2/"))


def install_modules():
    err = subprocess.run(["yarn", "--cwd", react_path, "install"]).stderr
    if err is not None:
        print(err)


def build_react_yarn():
    err = subprocess.run(["yarn", "--cwd", react_path, "build"]).stderr
    if err is not None:
        print(err)


def move_build():
    source = os.path.join(react_path, "build")
    dest = os.path.abspath(os.path.join(react_path, "../dist"))

    if os.path.exists(dest):
        shutil.rmtree(dest)
    shutil.move(source, dest)


if __name__ == "__main__":
    install_modules()
    build_react_yarn()
    move_build()
