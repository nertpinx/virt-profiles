#
# This file is part of the virt-profiles project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Copyright 2018 Red Hat, Inc.
#

FROM fedora:28

MAINTAINER "Francesco Romani" <fromani@redhat.com>

RUN groupadd -g 1001 virtprofiles && \
	useradd -u 1001 --create-home -g virtprofiles -s /bin/bash virtprofiles && \
	mkdir /usr/share/virtprofiles && \
	chown virtprofiles:virtprofiles /usr/share/virtprofiles && \
	chmod 0755 /usr/share/virtprofiles
WORKDIR /home/virtprofiles
USER 1001

COPY virtprofilesd /usr/bin/virtprofilesd

ENTRYPOINT [ "/usr/bin/virtprofilesd" ]
