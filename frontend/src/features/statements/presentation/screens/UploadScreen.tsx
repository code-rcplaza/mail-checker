import { AppShell } from "../../../../components/layout/AppShell";
import { BottomNav } from "../../../../components/layout/BottomNav";
import { BrandHeader } from "../../../../components/layout/BrandHeader";
import { SideNav } from "../../../../components/layout/SideNav";
import { useNavigateTo } from "../../../../lib/router/useNavigateTo";
import { useBanks } from "../../application/useBanks";
import { useUploadStatement } from "../../application/useUploadStatement";
import { SourceZone } from "../components/SourceZone";
import { UploadForm } from "../components/UploadForm";

export interface UploadScreenProps {
  /** The signed-in person, from the route guard. */
  readonly account: { readonly name: string; readonly email: string };
}

/** Elegir fuente: in the MVP, only uploading a statement. */
export function UploadScreen({ account }: UploadScreenProps) {
  const navigateTo = useNavigateTo();
  const banks = useBanks();
  const upload = useUploadStatement({ onUploaded: (result) => navigateTo(`/analisis/${result.id}`) });

  return (
    <AppShell
      header={<BrandHeader session={{ name: account.name, onNavigate: navigateTo }} />}
      zone={<SourceZone onNavigate={navigateTo} />}
      nav={<BottomNav onNavigate={navigateTo} />}
      sideNav={<SideNav account={account} onNavigate={navigateTo} />}
    >
      <UploadForm
        file={upload.file}
        fileError={upload.fileError}
        onPick={upload.pick}
        banks={banks}
        bank={upload.bank}
        onChooseBank={upload.chooseBank}
        needsBillingDate={upload.needsBillingDate}
        billingDate={upload.billingDateText}
        onWriteBillingDate={upload.writeBillingDate}
        error={upload.error}
        canUpload={upload.canUpload}
        isUploading={upload.isUploading}
        onUpload={upload.upload}
      />
    </AppShell>
  );
}
