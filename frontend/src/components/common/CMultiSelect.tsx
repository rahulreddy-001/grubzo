import { Text } from "@radix-ui/themes";
import Select, {
  type StylesConfig,
  type GroupBase,
  components,
} from "react-select";
import { X } from "lucide-react";

const styles: StylesConfig<OptionType, true, GroupBase<OptionType>> = {
  control: (base, state) => ({
    ...base,
    borderRadius: "4px",
    boxShadow: "none",
    padding: "0 0px",
    borderColor: state.isFocused ? "var(--accent-9)" : "var(--gray-6)",
    "&:hover": {
      borderColor: "var(--accent-9)",
    },
  }),

  menuPortal: (base) => ({
    ...base,
    zIndex: 999999,
    position: "absolute",
  }),

  menu: (base) => ({
    ...base,
    zIndex: 999999,
    position: "absolute",
    margin: 0,
    padding: 0,
  }),

  placeholder: (base) => ({
    ...base,
    fontSize: "13px",
    fontWeight: 400,
  }),

  option: (base) => ({
    ...base,
    fontSize: "12px",
    cursor: "pointer",
    padding: "6px 10px",
    margin: 0,
  }),

  multiValue: (base) => ({
    ...base,
    borderRadius: "3px",
    padding: "2px 4px",
    fontSize: "13px",
  }),
  dropdownIndicator: () => ({ display: "none" }),
  indicatorSeparator: () => ({ display: "none" }),
};

type OptionType = {
  value: string;
  label: string;
};

interface Props {
  selected: string[];
  options: string[];
  placeholder: string;
  onChange: (values: string[]) => void;
}
export default function CMultiSelect({
  selected,
  options,
  placeholder,
  onChange,
}: Props) {
  const optionList = options.map((o) => ({
    value: o,
    label: o.toLocaleUpperCase(),
  }));
  const selectedOptions = optionList.filter((o) => selected.includes(o.value));

  return (
    <Select<OptionType, true, GroupBase<OptionType>>
      isMulti
      options={optionList}
      value={selectedOptions}
      onChange={(vals) => onChange(vals.map((v) => v.value))}
      placeholder={<Text>{placeholder}</Text>}
      closeMenuOnSelect={false}
      isClearable={false}
      styles={styles}
      menuPortalTarget={document.body}
      menuPosition="absolute"
      components={{
        Option: (props) => (
          <components.Option {...props}>
            <Text
              size="1"
              style={{
                fontFamily:
                  "var(--font-system, system-ui, -apple-system, 'Segoe UI', Roboto)",
              }}
            >
              {props.label}
            </Text>
          </components.Option>
        ),
        MultiValueRemove: (props) => (
          <div
            {...props.innerProps}
            style={{
              display: "flex",
              alignItems: "center",
              cursor: "pointer",
              paddingLeft: "2px",
            }}
          >
            <X size={12} strokeWidth={2} />
          </div>
        ),
      }}
    />
  );
}
