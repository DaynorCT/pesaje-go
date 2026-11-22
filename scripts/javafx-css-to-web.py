#!/usr/bin/env python3
"""Rough JavaFX CSS to web CSS for T1 scaffold."""
import re
import pathlib
import sys

MAPPING = {
    "background-color": "background-color",
    "text-fill": "color",
    "font-family": "font-family",
    "font-size": "font-size",
    "font-weight": "font-weight",
    "border-color": "border-color",
    "border-radius": "border-radius",
    "border-width": "border-width",
    "padding": "padding",
    "background-radius": "border-radius",
    "opacity": "opacity",
    "cursor": "cursor",
}

SKIP = {"effect", "background-insets", "min-height", "max-height", "pref-height", "min-width", "max-width", "pref-width"}


def pad_val(val: str) -> str:
    if "px" in val:
        return val
    parts = val.split()
    if len(parts) == 4:
        return " ".join(p + ("px" if not p.endswith("px") else "") for p in parts)
    if len(parts) == 2:
        return f"{parts[0]}px {parts[1]}px"
    if len(parts) == 1 and parts[0].replace(".", "").isdigit():
        return f"{parts[0]}px"
    return val


def convert(src: pathlib.Path, dst: pathlib.Path) -> None:
    text = src.read_text(encoding="utf-8")
    rules = []
    for selector, body in re.findall(r"([^{]+)\{([^}]+)\}", text, re.S):
        selector = selector.strip()
        if not selector or selector.startswith("/*"):
            continue
        props = []
        for chunk in body.split(";"):
            chunk = chunk.strip()
            if not chunk.startswith("-fx-"):
                continue
            chunk = chunk[4:]
            if ":" not in chunk:
                continue
            prop, val = chunk.split(":", 1)
            prop, val = prop.strip(), val.strip()
            if prop in SKIP:
                continue
            css = MAPPING.get(prop, prop)
            if prop == "padding":
                val = pad_val(val)
            props.append(f"    {css}: {val};")
        if props:
            rules.append(selector + " {\n" + "\n".join(props) + "\n}")
    header = f"/* Web port of {src.name} (T1). Source: jacobitus-escritorio/pesaje/src/main/resources/{src.name} */\n\n"
    dst.write_text(header + "\n\n".join(rules) + "\n", encoding="utf-8")


if __name__ == "__main__":
    base = pathlib.Path(sys.argv[1])
    out = pathlib.Path(sys.argv[2])
    out.mkdir(parents=True, exist_ok=True)
    convert(base / "pesaje.css", out / "pesaje.css")
    convert(base / "jacobitus.css", out / "jacobitus.css")
