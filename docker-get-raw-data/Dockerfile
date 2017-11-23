# Copyright (C) Activision Publishing, Inc. 2017
# https://github.com/Demonware/harbor-analytics
# Author: David Rieger
# Licensed under the 3-Clause BSD License (the "License");
# you may not use this file except in compliance with the License.

FROM garland/aws-cli-docker
ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
