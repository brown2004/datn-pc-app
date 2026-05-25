type ToggleSwitchProps = {
  checked: boolean;
  disabled?: boolean;
  onChange: (checked: boolean) => void;
};

export function ToggleSwitch({ checked, disabled, onChange }: ToggleSwitchProps) {
  return (
    <button
      className={`relative h-[30px] w-[54px] cursor-pointer rounded-full border transition-colors disabled:cursor-wait disabled:opacity-60 ${
        checked ? "border-[#0f766e] bg-[#0f766e]" : "border-[#c7ced8] bg-[#e5e7eb]"
      }`}
      type="button"
      disabled={disabled}
      onClick={() => onChange(!checked)}
      aria-pressed={checked}
    >
      <span
        className={`absolute left-[3px] top-[3px] h-[22px] w-[22px] rounded-full bg-white shadow-[0_1px_3px_rgba(15,23,42,0.22)] transition-transform ${
          checked ? "translate-x-6" : "translate-x-0"
        }`}
      />
    </button>
  );
}
