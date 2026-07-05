import { useState } from "react";

/**
 * StarRating — hover vs selected state, derived rendering.
 *
 * Props:
 *   count         number    how many stars (default 5)
 *   defaultValue  number    initially selected rating (default 0)
 *   onChange      function  called with the new rating when a star is clicked
 *
 * Required UI (the tests rely on this):
 *   - `count` buttons, each with aria-label "Rate N star" / "Rate N stars"
 *   - each button has attribute data-filled="true" or "false"
 *   - a star is FILLED when its position <= the value currently shown
 *
 * Rules:
 *   - The value shown is the HOVERED rating if the user is hovering a star,
 *     otherwise the SELECTED rating.
 *   - Hovering star N fills stars 1..N. Moving the mouse away restores the
 *     selected rating.
 *   - Clicking star N selects N and calls onChange(N).
 */
export default function StarRating({ count = 5, defaultValue = 0, onChange }) {
  const [selected, setSelected] = useState(defaultValue);
  const [hovered, setHovered] = useState(null);

  const displayValue = hovered ?? selected;

  const select = (n) => {
    setSelected(n);
    onChange?.(n);
  };

  return (
    <div>
      {Array.from({ length: count }, (_, i) => {
        const position = i + 1;
        const filled = position <= displayValue;
        return (
          <button
            key={position}
            aria-label={`Rate ${position} star${position === 1 ? "" : "s"}`}
            data-filled={filled}
            onClick={() => select(position)}
            onMouseEnter={() => {
              setHovered(position);
            }}
            onMouseLeave={() => {
              setHovered(null);
            }}
          >
            {filled ? "★" : "☆"}
          </button>
        );
      })}
    </div>
  );
}
