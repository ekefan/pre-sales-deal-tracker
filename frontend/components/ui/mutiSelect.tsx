import * as React from "react";
import { Checkbox } from "@/components/ui/checkbox"; // Adjust the import according to your project
import { Button } from "@/components/ui/button"; // Adjust the import according to your project
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from "@/components/ui/popover";
import {
  Command,
  CommandInput,
  CommandList,
  CommandGroup,
  CommandItem,
} from "@/components/ui/command";
import { XCircle, ChevronDown } from "lucide-react";
import { cn } from "@/lib/utils"; // Adjust the import according to your project
import { Badge } from "@/components/ui/badge"; // Import the Badge component

interface MultiSelectOption {
  label: string;
  value: string;
}

interface MultiSelectProps {
  options: MultiSelectOption[];
  value: string[];
  onChange: (value: string[]) => void;
  placeholder?: string;
  className?: string;
}

export const MultiSelect: React.FC<MultiSelectProps> = ({
  options,
  value,
  onChange,
  placeholder = "Select services",
  className,
}) => {
  const [isOpen, setIsOpen] = React.useState(false);

  const toggleOption = (optionValue: string) => {
    const newValues = value.includes(optionValue)
      ? value.filter((val) => val !== optionValue)
      : [...value, optionValue];
    onChange(newValues);
  };

  return (
    <Popover open={isOpen} onOpenChange={setIsOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          className={cn("w-full justify-between flex-grow h-full text-left text-sm  font-normal text-slate-600", className)}
        >
         <div className="flex flex-wrap gap-1">
         {value.length > 0
            ? value.map((val) => {
                const option = options.find((opt) => opt.value === val);
                return option ? (
                  <Badge key={val} className="mr-2">
                    {option.label}
                    <XCircle
                      className="ml-2 h-4 w-4 cursor-pointer"
                      onClick={(e) => {
                        e.stopPropagation();
                        toggleOption(option.value);
                      }}
                    />
                  </Badge>
                ) : null;
              })
            : placeholder}
         </div>
          <ChevronDown className="ml-2"  absoluteStrokeWidth={true}/>
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-full p-0">
        <Command>
          <CommandInput placeholder="Search..." />
          <CommandList>
            <CommandGroup>
              {options.map((option) => (
                <CommandItem
                  key={option.value}
                  onSelect={() => toggleOption(option.value)}
                >
                  <Checkbox
                    checked={value.includes(option.value)}
                    onCheckedChange={() => toggleOption(option.value)}
                  />
                  {option.label}
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
};