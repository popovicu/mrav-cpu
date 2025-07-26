def generate_bundling_shell(files, output_file):
    # Create the command that concatenates all sources with newlines between them
    cmd = "; printf \"\\n\\n\"; ".join(["cat " + f.path for f in files])

    # Add redirection to output file
    return "(" + cmd + ") > " + output_file.path
