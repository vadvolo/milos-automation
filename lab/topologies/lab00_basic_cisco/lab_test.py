import os
import subprocess
import unittest


class TestLab(unittest.TestCase):
    def test_gen(self):
        # exec -u root -t -i annet annet gen lab-r1.nh.com lab-r2.nh.com lab-r3.nh.com
        genInfo = os.popen(
            "docker exec -u root -t -i annet annet gen lab-r1.nh.com lab-r2.nh.com lab-r3.nh.com",
        ).read()
        self.assertEqual(
            genInfo,
            "\x1b[1m\x1b[42m"
            + "# -------------------- lab-r1.nh.com.cfg --------------------\x1b[49m\x1b[0m\ninterface FastEthernet0/0\n  description disconnected\n  mtu 1500\ninterface FastEthernet0/1\n  description disconnected\n  mtu 1500\ninterface GigabitEthernet1/0\n  description to_lab-r2.nh.com_GigabitEthernet1/0\n  mtu 4000\ninterface GigabitEthernet2/0\n  description disconnected\n  mtu 1500\n"
            + "\x1b[1m\x1b[42m# -------------------- lab-r2.nh.com.cfg --------------------\x1b[49m\x1b[0m\ninterface FastEthernet0/0\n  description disconnected\n  mtu 1500\ninterface FastEthernet0/1\n  description disconnected\n  mtu 1500\ninterface GigabitEthernet1/0\n  description to_lab-r1.nh.com_GigabitEthernet1/0\n  mtu 4000\ninterface GigabitEthernet2/0\n  description to_lab-r3.nh.com_GigabitEthernet1/0\n  mtu 4000\n"
            + "\x1b[1m\x1b[42m# -------------------- lab-r3.nh.com.cfg --------------------\x1b[49m\x1b[0m\ninterface FastEthernet0/0\n  description disconnected\n  mtu 1500\ninterface FastEthernet0/1\n  description disconnected\n  mtu 1500\ninterface GigabitEthernet1/0\n  description to_lab-r2.nh.com_GigabitEthernet2/0\n  mtu 4000\ninterface GigabitEthernet2/0\n  description disconnected\n  mtu 1500\n\x1b[0m",
        )

    def test_diff(self):
        # exec -u root -t -i annet annet diff lab-r1.nh.com lab-r2.nh.com lab-r3.nh.com
        diffInfo = os.popen(
            "docker exec -u root -t -i annet annet diff lab-r1.nh.com lab-r2.nh.com lab-r3.nh.com"
        ).read()
        self.assertEqual(
            diffInfo,
            "\x1b[1m\x1b[42m"
            + "# -------------------- lab-r2.nh.com.cfg --------------------\x1b[49m\x1b[0m\n\x1b[1m\x1b[36m  interface FastEthernet0/0\x1b[0m\n\x1b[1m\x1b[32m+   description disconnected\x1b[0m\n\x1b[1m\x1b[36m  interface FastEthernet0/1\x1b[0m\n\x1b[1m\x1b[32m+   description disconnected\x1b[0m\n\x1b[1m\x1b[36m  interface GigabitEthernet1/0\x1b[0m\n\x1b[1m\x1b[32m+   description to_lab-r1.nh.com_GigabitEthernet1/0\x1b[0m\n\x1b[1m\x1b[32m+   mtu 4000\x1b[0m\n\x1b[1m\x1b[36m  interface GigabitEthernet2/0\x1b[0m\n\x1b[1m\x1b[32m+   description to_lab-r3.nh.com_GigabitEthernet1/0\x1b[0m\n\x1b[1m\x1b[32m+   mtu 4000\x1b[0m\n\x1b[1m\x1b[42m"
            + "# -------------------- lab-r1.nh.com.cfg --------------------\x1b[49m\x1b[0m\n\x1b[1m\x1b[36m  interface FastEthernet0/0\x1b[0m\n\x1b[1m\x1b[32m+   description disconnected\x1b[0m\n\x1b[1m\x1b[36m  interface FastEthernet0/1\x1b[0m\n\x1b[1m\x1b[32m+   description disconnected\x1b[0m\n\x1b[1m\x1b[36m  interface GigabitEthernet1/0\x1b[0m\n\x1b[1m\x1b[32m+   description to_lab-r2.nh.com_GigabitEthernet1/0\x1b[0m\n\x1b[1m\x1b[32m+   mtu 4000\x1b[0m\n\x1b[1m\x1b[36m  interface GigabitEthernet2/0\x1b[0m\n\x1b[1m\x1b[32m+   description disconnected\x1b[0m\n\x1b[1m\x1b[42m"
            + "# -------------------- lab-r3.nh.com.cfg --------------------\x1b[49m\x1b[0m\n\x1b[1m\x1b[36m  interface FastEthernet0/0\x1b[0m\n\x1b[1m\x1b[32m+   description disconnected\x1b[0m\n\x1b[1m\x1b[36m  interface FastEthernet0/1\x1b[0m\n\x1b[1m\x1b[32m+   description disconnected\x1b[0m\n\x1b[1m\x1b[36m  interface GigabitEthernet1/0\x1b[0m\n\x1b[1m\x1b[32m+   description to_lab-r2.nh.com_GigabitEthernet2/0\x1b[0m\n\x1b[1m\x1b[32m+   mtu 4000\x1b[0m\n\x1b[1m\x1b[36m  interface GigabitEthernet2/0\x1b[0m\n\x1b[1m\x1b[32m+   description disconnected\x1b[0m\n\x1b[0m"
        )

    def test_patch(self):
        # exec -u root -t -i annet annet patch lab-r1.nh.com lab-r2.nh.com lab-r3.nh.com
        patchInfo = os.popen(
            "docker exec -u root -t -i annet annet patch lab-r1.nh.com lab-r2.nh.com lab-r3.nh.com"
        ).read()
        self.assertEqual(
            patchInfo,
            "\x1b[1m\x1b[42m"
            + "# -------------------- lab-r1.nh.com.patch --------------------\x1b[49m\x1b[0m\ninterface FastEthernet0/0\n  description disconnected\n  exit\ninterface FastEthernet0/1\n  description disconnected\n  exit\ninterface GigabitEthernet1/0\n  description to_lab-r2.nh.com_GigabitEthernet1/0\n  mtu 4000\n  exit\ninterface GigabitEthernet2/0\n  description disconnected\n  exit\n\x1b[1m\x1b[42m"
            + "# -------------------- lab-r2.nh.com.patch --------------------\x1b[49m\x1b[0m\ninterface FastEthernet0/0\n  description disconnected\n  exit\ninterface FastEthernet0/1\n  description disconnected\n  exit\ninterface GigabitEthernet1/0\n  description to_lab-r1.nh.com_GigabitEthernet1/0\n  mtu 4000\n  exit\ninterface GigabitEthernet2/0\n  description to_lab-r3.nh.com_GigabitEthernet1/0\n  mtu 4000\n  exit\n\x1b[1m\x1b[42m"
            + "# -------------------- lab-r3.nh.com.patch --------------------\x1b[49m\x1b[0m\ninterface FastEthernet0/0\n  description disconnected\n  exit\ninterface FastEthernet0/1\n  description disconnected\n  exit\ninterface GigabitEthernet1/0\n  description to_lab-r2.nh.com_GigabitEthernet2/0\n  mtu 4000\n  exit\ninterface GigabitEthernet2/0\n  description disconnected\n  exit\n\x1b[0m"
        )


if __name__ == "__main__":
    unittest.main(verbosity=2)

# How to get output from stdout
# genInfo = subprocess.run(
#     [
#         "docker",
#         "exec",
#         "-u",
#         "root",
#         "-t",
#         "-i",
#         "annet",
#         "annet",
#         "gen",
#         "lab-r1.nh.com",
#         "lab-r2.nh.com",
#         "lab-r3.nh.com",
#     ],
#     capture_output=True,
#     text=True,
# )
