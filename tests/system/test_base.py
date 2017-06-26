from testpool-beat import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting Testpool-beat normally
        """
        self.render_config_template(
                path=os.path.abspath(self.working_dir) + "/log/*"
        )

        testpool-beat_proc = self.start_beat()
        self.wait_until( lambda: self.log_contains("testpool-beat is running"))
        exit_code = testpool-beat_proc.kill_and_wait()
        assert exit_code == 0
