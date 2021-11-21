import unittest
import testinfra
import sys
from os import path
sys.path.append( path.dirname( path.dirname( path.abspath(__file__) ) ) )
import app


class TestINFRA(unittest.TestCase):
    def setUp(self):
        self.host = testinfra.get_host("local://", sudo=False)

    def test_lirc_package(self):
        lirc = self.host.package("lirc")
        self.assertTrue(lirc.is_installed)

    def test_nginx_package(self):
        nginx = self.host.package("nginx")
        self.assertTrue(nginx.is_installed)

if __name__ == '__main__':
    unittest.main()
