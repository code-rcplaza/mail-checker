import { AppShell } from "../../../../components/layout/AppShell";
import { BottomNav } from "../../../../components/layout/BottomNav";
import { BrandHeader } from "../../../../components/layout/BrandHeader";
import { SideNav } from "../../../../components/layout/SideNav";
import { useNavigateTo } from "../../../../lib/router/useNavigateTo";
import { useAnalysis } from "../../application/useAnalysis";
import { AnalysisResult, AnalysisZone } from "../components/AnalysisResult";
import { UploadLoadingContent, UploadLoadingZone } from "../components/UploadLoading";

export interface AnalysisScreenProps {
  readonly statementId: string;
  readonly account: { readonly name: string; readonly email: string };
}

/** Análisis: what was read and how many subscriptions came out of it. */
export function AnalysisScreen({ statementId, account }: AnalysisScreenProps) {
  const navigateTo = useNavigateTo();
  const analysis = useAnalysis(statementId);
  const { found } = analysis;

  return (
    <AppShell
      header={<BrandHeader session={{ name: account.name, onNavigate: navigateTo }} />}
      zone={found === undefined ? <UploadLoadingZone /> : <AnalysisZone found={found} />}
      nav={<BottomNav onNavigate={navigateTo} />}
      sideNav={<SideNav account={account} onNavigate={navigateTo} />}
    >
      {analysis.error !== undefined && (
        <p role="alert" className="mx-6 mt-8 rounded-md bg-error-surface p-4 text-sm text-danger-ink md:mx-0">
          {analysis.error}
        </p>
      )}
      {found === undefined && analysis.isLoading && <UploadLoadingContent />}
      {found !== undefined && (
        <AnalysisResult found={found} uploaded={analysis.uploaded} onNavigate={navigateTo} />
      )}
    </AppShell>
  );
}
