import { useQuery } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { AuthGateway } from "../../features/auth/application/ports";
import { useGateways } from "./GatewaysProvider";
import { renderWithGateways } from "./renderWithGateways";

const ana = { id: "u1", name: "Ana Prueba", email: "ana@example.com", initials: "AP" };

function MeName() {
  const { auth } = useGateways();
  const me = useQuery({ queryKey: ["me"], queryFn: () => auth.getMe() });
  return <p>{me.data?.name ?? "…"}</p>;
}

describe("gateways", () => {
  it("hands the hooks the double the provider got", async () => {
    const auth: AuthGateway = { getMe: vi.fn(async () => ana) };

    renderWithGateways(<MeName />, { auth });

    expect(await screen.findByText("Ana Prueba")).toBeInTheDocument();
    expect(auth.getMe).toHaveBeenCalledOnce();
  });

  it("fails loudly without a provider", () => {
    vi.spyOn(console, "error").mockImplementation(() => {});

    expect(() => render(<MeName />)).toThrow("useGateways needs a GatewaysProvider above it");
  });

  it("gives each test a new query client", () => {
    const first = renderWithGateways(<p>uno</p>);
    const second = renderWithGateways(<p>dos</p>);

    expect(first.queryClient).not.toBe(second.queryClient);
  });

  it("fails a test that uses a gateway it did not give", async () => {
    renderWithGateways(<MeName />);

    expect(await screen.findByText("…")).toBeInTheDocument();
  });
});
