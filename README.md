# Photobook

It's a simple Single Page Application that allows to upload photos and see them in the homepage.
I created it to have a place to share images for a specific event/party.  
The upload page is without access control, so everyone can upload photos.

## Features

- No login, I wanted to keep it simple as possible;
- The homepage is instantaneously updated when a new photo is uploaded.

## Develop

### Requirements

To develop the webservice, you need:
- Vagrant with NFS support (to share the foldersrc/)
- libvirt

You can follow these guides for Vagrant/libvirt:  
https://linuxsimba.com/vagrant-libvirt-install  
https://docs.cumulusnetworks.com/display/VX/Vagrant+and+Libvirt+with+KVM+or+QEMU  
https://docs.cumulusnetworks.com/cumulus-vx/Getting-Started/Libvirt-and-KVM-QEMU/

In summary the steps to do are:
```
sudo apt install vagrant vagrant-libvirt
sudo apt install libvirt-daemon-system  python3-distutils  python3-gi-cairo  python3-lib2to3  python3-libvirt qemu-kvm   spice-client-glib-usb-acl-helper  systemd-container    virt-manager  virt-viewer   virtinst  libgovirt-common   gir1.2-gtk-vnc-2.0
sudo systemctl restart libvirtd
# logout and login, you should have the libvirt group
sudo apt install python3-paramiko python3-venv python3-pip
```

### Steps to develop
```
user@host:~ cd photobook/
user@host:~/photobook vagrant-tools/dev-web-back.sh
# another shell
user@host:~ cd photobook/
user@host:~/photobook vagrant-tools/dev-web-front.sh
```

## Deployment

The following steps will create a local stack in Docker. The script creates from the source the new images and start all the necessary components to serve the single page application.
```
user@host:~ cd photobook/
user@host:~/photobook docker/lib/createLocalStack.sh
```
