import { ButtonHTMLAttributes, forwardRef } from "react";
import { cn } from "../../lib/utils";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "outline" | "ghost" | "danger";
  size?: "sm" | "md" | "lg" | "icon";
  isLoading?: boolean;
}

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = "primary", size = "md", isLoading, children, disabled, ...props }, ref) => {
    const variants = {
      primary: "bg-blue-600 text-white shadow-lg shadow-blue-100 hover:bg-blue-500 active:scale-[0.98]",
      secondary: "bg-gray-100 text-gray-900 hover:bg-gray-200 active:scale-[0.98]",
      outline: "border border-gray-200 bg-white text-gray-600 hover:bg-gray-50 active:scale-[0.98]",
      ghost: "text-gray-500 hover:bg-gray-100 hover:text-gray-900",
      danger: "bg-red-600 text-white shadow-lg shadow-red-100 hover:bg-red-500 active:scale-[0.98]",
    };

    const sizes = {
      sm: "px-3 py-1.5 text-xs font-semibold",
      md: "px-6 py-2.5 text-sm font-bold",
      lg: "px-8 py-3 text-base font-bold",
      icon: "p-2",
    };

    return (
      <button
        ref={ref}
        disabled={disabled || isLoading}
        className={cn(
          "inline-flex items-center justify-center gap-2 rounded-xl transition-all disabled:opacity-50",
          variants[variant],
          sizes[size],
          className
        )}
        {...props}
      >
        {children}
      </button>
    );
  }
);

Button.displayName = "Button";
export { Button };
