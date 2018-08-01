The VM configuration pipeline
=============================

Concepts
--------

The VM configuration pipeline is the process which takes an input a VirtualMachineInstance definition, and produces as output a libvirt Domain Spec.

This process is called "pipeline" because is composed by few stages; to describe them, we need to introduce first how the abstract concept of VM configuration is layered.

KubeVirt depends on libvirt to manage virtual machines; so, at the bottom, we have the libvirt layer. Libvirt describes VMs using the Domain Spec, in XML format.
Let's call this layer LDS (Libvirt Domain Spec). This layer fully specifies what resources the VMs exposed to the guest, and how those virtual resources are mapped into host resources.

The emulated configuration presented to the guest is called the Frontend, while the mapping (and the tuning of) between Virtualized resources and Host resources is called the Backend.

Users knows and care about the Frontend of the LDS. It is up to the administator, or to the management application, to figure out the optimal Backend.
It is worth nothing that in the LDS the Frontend and Backend section are deeply intertwined in a single document.

Since users want to reason in term of the requests of the VM, they may not want to care to specify all the fine details of the Frontend spec, and the may find useful a even more
abstract specification. This is what the Kubevirt's VirtualMachineInstance Definition offers (KVMID), an abstract VM definition that allows the users to focus on the kye aspect they care about,
leaving the system to translate the KVMID in the LDS frontend.


Stages
------

* Stage 1:
receive KVMID, augment with presets, accept into the system

```
[user] -{KVMID.yaml}->[STAGE#1]->-{KVMID.yaml}
                         A
                         |
        -{preset_1.yaml}-+
        -{preset_2.yaml}-+
          ...            |
        -{preset_N.yaml}-'
```

This is what is already implemented in kubevirt, The KVMID is augmented with the data from presets. The KVMID is subset of the LDS Frontend, so the configuration is not complete yet.

* Stage 2:
Translate the KVMID in a fully-specified LDS Frontend

```
{KVMID.yaml}->[STAGE#2]->-{Yet Undefined Format}
                 A
                 |
    -{profiles?>-+
```

A fully specified LDS Frontend has not a clear representation. There are three major options:
1. represent the LDS Frontend with a partial LDS, to be filled up later. Possible today.
2. represent the LDS Frontend with an enhanced KVMID, yet to be done
3. represent the LDS Frontend with another intermediate format. Not even planned. Mentione for the sake of completeness, but I believe either #1 or #2 are stricly better alternatives.


* Stage 3:
Sort out the LDS Backend, and complete the LDS

```
{Yet Undefined Format}->[STAGE#3]->-{LDS.xml}
                           A
                           |
              -{profiles?>-+
```


The system figures out the optimal LDS Backend; now the system has all the information to emit a complete LDS, which is the final result of the pipeline.

The current implementation in KubeVirt
--------------------------------------

KubeVirt currently implements the pipeline above, even though not in a clean fashion - their implementation predates this document -, and collapses stage #2 and #3 in a single stage.
The single combined KubeVirt #2+#3 stage translates the KVMID in a LDS-equivalent format *and* perform changes to determine the LDS backend. Later, it trivially marshals the LDS-equivalent
format in a LDS XML.


Future implementation
---------------------

TBD

key ideas:
1. clone and extend profiles; apply them in stage #1
2. expect KVMID to catch up with LDS Frontend, use extended Presets to tune LDS Frontend
3. use partial LDS as Yet Undefined Format in stage #2/#3
4. require a different set of profiles to tune the LDS Backend.
