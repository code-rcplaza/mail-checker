import { useId } from "react";

export interface SegmentedOption {
  readonly value: string;
  readonly label: string;
}

export interface SegmentedControlProps {
  readonly label: string;
  readonly options: ReadonlyArray<SegmentedOption>;
  readonly value: string;
  readonly onChange: (value: string) => void;
}

/**
 * Native radios drawn as segments: the browser gives one tab stop, arrows that
 * move and choose, and the checked state for screen readers.
 */
export function SegmentedControl({ label, options, value, onChange }: SegmentedControlProps) {
  const name = useId();
  return (
    <fieldset className="inline-flex rounded-md bg-field p-1">
      <legend className="sr-only">{label}</legend>
      {options.map((option) => (
        <label
          key={option.value}
          className="flex min-h-11 flex-1 cursor-pointer items-center justify-center rounded-sm px-4 text-sm font-medium text-ink-secondary has-checked:bg-page has-checked:text-ink has-focus-visible:outline-2 has-focus-visible:outline-primary"
        >
          <input
            type="radio"
            name={name}
            value={option.value}
            checked={option.value === value}
            onChange={() => onChange(option.value)}
            className="sr-only"
          />
          {option.label}
        </label>
      ))}
    </fieldset>
  );
}
