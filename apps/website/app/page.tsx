import Navbar from "@/components/Navbar";
import Hero from "@/components/Hero";
import Terminal from "@/components/Terminal";
import Categories from "@/components/Categories";
import Profiles from "@/components/Profiles";
import GithubScanner from "@/components/GithubScanner";
import Skills from "@/components/Skills";
import InstallMethods from "@/components/InstallMethods";
import Footer from "@/components/Footer";
import UpdatePopup from "@/components/UpdatePopup";
import SectionWrapper from "@/components/SectionWrapper";
import SectionDivider from "@/components/SectionDivider";
import HowItWorks from "@/components/HowItWorks";
import RealExample from "@/components/RealExample";
import AudienceAndComparison from "@/components/AudienceAndComparison";
import Testimonials from "@/components/Testimonials";
import CommandCenterPreview from "@/components/CommandCenterPreview";

export default function Home() {
  return (
    <main className="min-h-screen bg-grid">
      <Navbar />
      <Hero />
      <CommandCenterPreview />

      <SectionDivider variant="wave" />
      <SectionWrapper gradient>
        <HowItWorks />
      </SectionWrapper>

      <SectionDivider variant="angle" flip />
      <SectionWrapper gradient>
        <RealExample />
      </SectionWrapper>

      <SectionDivider variant="curve" />
      <SectionWrapper gradient>
        <AudienceAndComparison />
      </SectionWrapper>

      <SectionDivider variant="wave" flip />
      <SectionWrapper gradient>
        <Terminal />
      </SectionWrapper>

      <SectionDivider variant="angle" />
      <SectionWrapper gradient>
        <Categories />
      </SectionWrapper>

      <SectionDivider variant="curve" flip />
      <SectionWrapper gradient>
        <Profiles />
      </SectionWrapper>

      <SectionDivider variant="wave" />
      <SectionWrapper gradient>
        <GithubScanner />
      </SectionWrapper>

      <SectionDivider variant="angle" flip />
      <SectionWrapper gradient>
        <Skills />
      </SectionWrapper>

      <SectionDivider variant="curve" />
      <SectionWrapper gradient>
        <Testimonials />
      </SectionWrapper>

      <SectionDivider variant="wave" flip />
      <SectionWrapper gradient>
        <InstallMethods />
      </SectionWrapper>

      <Footer />
      <UpdatePopup />
    </main>
  );
}
