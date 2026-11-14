import { fireEvent, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";
import { BankSelect } from "./BankSelect";
import { UploadDropzone } from "./UploadDropzone";

const pdf = new File(["%PDF"], "cartola-octubre.pdf", { type: "application/pdf" });

describe("UploadDropzone", () => {
  afterEach(() => vi.restoreAllMocks());

  it("takes a chosen PDF and shows its name", async () => {
    const onPick = vi.fn();
    const { rerender } = render(<UploadDropzone file={undefined} error={undefined} onPick={onPick} />);

    await userEvent.upload(screen.getByLabelText("Archivo de la cartola"), pdf);
    expect(onPick).toHaveBeenCalledWith(pdf);

    rerender(<UploadDropzone file={pdf} error={undefined} onPick={onPick} />);
    expect(screen.getByText("cartola-octubre.pdf")).toBeInTheDocument();
  });

  it("takes a dropped PDF", () => {
    const onPick = vi.fn();
    render(<UploadDropzone file={undefined} error={undefined} onPick={onPick} />);

    fireEvent.drop(screen.getByTestId("dropzone"), { dataTransfer: { files: [pdf] } });

    expect(onPick).toHaveBeenCalledWith(pdf);
  });

  it.each(["{Enter}", " "])("opens the file picker with %j from the keyboard", async (key) => {
    const openPicker = vi.spyOn(HTMLInputElement.prototype, "click").mockImplementation(() => {});
    render(<UploadDropzone file={undefined} error={undefined} onPick={() => {}} />);

    screen.getByRole("button", { name: /Elegir archivo PDF/ }).focus();
    await userEvent.keyboard(key);

    expect(openPicker).toHaveBeenCalled();
  });

  it("shows why a file was rejected", () => {
    render(<UploadDropzone file={undefined} error="Sube la cartola en PDF." onPick={() => {}} />);

    expect(screen.getByRole("alert")).toHaveTextContent("Sube la cartola en PDF.");
  });
});

describe("BankSelect", () => {
  it("offers automatic detection and each bank by its name", async () => {
    const onChange = vi.fn();
    render(<BankSelect banks={[{ id: "bci" }, { id: "nuevo" }]} value="" onChange={onChange} />);

    const select = screen.getByLabelText("Banco");
    expect([...select.querySelectorAll("option")].map((option) => option.textContent)).toEqual([
      "Detectar automáticamente",
      "BCI",
      "NUEVO",
    ]);
    await userEvent.selectOptions(select, "bci");
    expect(onChange).toHaveBeenCalledWith("bci");
  });
});
