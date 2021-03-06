package seqio

import (
	"strings"
	"testing"

	"github.com/go-gts/gts"
	"github.com/go-pars/pars"
)

var qualifierTests = []string{
	"                     /allele=\"adh1-1\"",
	"                     /altitude=\"-256 m\"",
	"                     /altitude=\"330.12 m\"",
	"                     /anticodon=(pos:34..36,aa:Phe,seq:aaa)",
	"                     /anticodon=(pos:join(5,495..496),aa:Leu,seq:taa)",
	"                     /anticodon=(pos:complement(4156..4158),\n                     aa:Gln,seq:ttg)",
	"                     /artificial_location=\"heterogeneous population sequenced\"",
	"                     /artificial_location=\"low-quality sequence region\"",
	"                     /bio_material=\"CGC:CB3912\"",
	"                     /bound_moiety=\"GAL4\"",
	"                     /cell_line=\"MCF7\"",
	"                     /cell_type=\"leukocyte\"",
	"                     /chromosome=\"1\"",
	"                     /citation=[3]",
	"                     /clone=\"lambda-hIL7.3\"",
	"                     /clone_lib=\"lambda-hIL7\"",
	"                     /codon_start=2",
	"                     /collected_by=\"Dan Janzen\"",
	"                     /collection_date=\"21-Oct-1952\"",
	"                     /collection_date=\"Oct-1952\"",
	"                     /collection_date=\"1952\"",
	"                     /collection_date=\"1952-10-21T11:43Z\"",
	"                     /collection_date=\"1952-10-21T11Z\"",
	"                     /collection_date=\"1952-10-21\"",
	"                     /collection_date=\"1952-10\"",
	"                     /collection_date=\"21-Oct-1952/15-Feb-1953\"",
	"                     /collection_date=\"Oct-1952/Feb-1953\"",
	"                     /collection_date=\"1952/1953\"",
	"                     /collection_date=\"1952-10-21/1953-02-15\"",
	"                     /collection_date=\"1952-10/1953-02\"",
	"                     /collection_date=\"1952-10-21T11:43Z/1952-10-21T17:43Z\"",
	"                     /collection_date=\"2015-10-11T17:53:03Z\"",
	"                     /compare=AJ634337.1",
	"                     /country=\"Canada:Vancouver\"",
	"                     /country=\"France:Cote d'Azur, Antibes\"",
	"                     /country=\"Atlantic Ocean:Charlie Gibbs Fracture Zone\"",
	"                     /cultivar=\"Nipponbare\"",
	"                     /cultivar=\"Tenuifolius\"",
	"                     /cultivar=\"Candy Cane\"",
	"                     /cultivar=\"IR36\"",
	"                     /culture_collection=\"ATCC:26370\"",
	"                     /db_xref=\"UniProtKB/Swiss-Prot:P28763\"",
	"                     /dev_stage=\"fourth instar larva\"",
	"                     /direction=LEFT",
	"                     /EC_number=\"1.1.2.4\"",
	"                     /EC_number=\"1.1.2.-\"",
	"                     /EC_number=\"1.1.2.n\"",
	"                     /EC_number=\"1.1.2.n1\"",
	"                     /ecotype=\"Columbia\"",
	"                     /environmental_sample",
	"                     /estimated_length=unknown",
	"                     /estimated_length=342",
	"                     /exception=\"RNA editing\"",
	"                     /exception=\"reasons given in citation\"",
	"                     /exception=\"rearrangement required for product\"",
	"                     /exception=\"annotated by transcript or proteomic data\"",
	"                     /experiment=\"5' RACE\"",
	"                     /experiment=\"Northern blot [DOI: 12.3456/FT.789.1.234-567.2010]\"",
	"                     /experiment=\"heterologous expression system of Xenopus laevis\n                     oocytes [PMID: 12345678, 10101010, 987654]\"",
	"                     /experiment=\"COORDINATES: 5' and 3' RACE\"",
	"                     /focus",
	"                     /frequency=\"23/108\"",
	"                     /frequency=\"1 in 12\"",
	"                     /frequency=\".85\"",
	"                     /function=\"essential for recognition of cofactor\"",
	"                     /gap_type=\"between scaffolds\"",
	"                     /gene=\"ilvE\"",
	"                     /gene_synonym=\"Hox-3.3\"",
	"                     /germline",
	"                     /haplogroup=\"H*\"",
	"                     /haplotype=\"Dw3 B5 Cw1 A1\"",
	"                     /host=\"Homo sapiens\"",
	"                     /host=\"Homo sapiens 12 year old girl\"",
	"                     /host=\"Rhizobium NGR234\"",
	"                     /identified_by=\"John Burns\"",
	"                     /inference=\"COORDINATES:profile:tRNAscan:2.1\"",
	"                     /inference=\"similar to DNA sequence:INSD:AY411252.1\"",
	"                     /inference=\"similar to RNA sequence, mRNA:RefSeq:NM_000041.2\"",
	"                     /inference=\"similar to DNA sequence (same\n                     species):INSD:AACN010222672.1\"",
	"                     /inference=\"protein motif:InterPro:IPR001900\"",
	"                     /inference=\"ab initio prediction:Genscan:2.0\"",
	"                     /inference=\"alignment:Splign:1.0\"",
	"                     /inference=\"alignment:Splign:1.26p:RefSeq:NM_000041.2,INSD:BC003557.1\"",
	"                     /isolate=\"Patient #152\"",
	"                     /isolate=\"DGGE band PSBAC-13\"",
	"                     /isolation_source=\"rumen isolates from standard\n                     Pelleted ration-fed steer #67\"",
	"                     /isolation_source=\"permanent Antarctic sea ice\"",
	"                     /isolation_source=\"denitrifying activated sludge from\n                     carbon_limited continuous reactor\"",
	"                     /lab_host=\"Gallus gallus\"",
	"                     /lab_host=\"Gallus gallus embryo\"",
	"                     /lab_host=\"Escherichia coli strain DH5 alpha\"",
	"                     /lab_host=\"Homo sapiens HeLa cells\"",
	"                     /lat_lon=\"47.94 N 28.12 W\"",
	"                     /lat_lon=\"45.0123 S 4.1234 E\"",
	"                     /linkage_evidence=\"paired-ends\"",
	"                     /linkage_evidence=\"within clone\"",
	"                     /locus_tag=\"ABC_0022\"",
	"                     /locus_tag=\"A1C_00001\"",
	"                     /macronuclear",
	"                     /mating_type=\"MAT-1\"",
	"                     /mating_type=\"plus\"",
	"                     /mating_type=\"-\"",
	"                     /mating_type=\"odd\"",
	"                     /mating_type=\"even\"",
	"                     /metagenome_source=\"human gut metagenome\"",
	"                     /metagenome_source=\"soil metagenome\"",
	"                     /mobile_element_type=\"transposon:Tnp9\"",
	"                     /mod_base=m5c",
	"                     /mol_type=\"genomic DNA\"",
	"                     /ncRNA_class=\"miRNA\"",
	"                     /ncRNA_class=\"siRNA\"",
	"                     /ncRNA_class=\"scRNA\"",
	"                     /note=\"This qualifier is equivalent to a comment.\"",
	"                     /number=4",
	"                     /number=6B",
	"                     /old_locus_tag=\"RSc0382\"",
	"                     /locus_tag=\"YPO0002\"",
	"                     /operon=\"lac\"",
	"                     /organelle=\"chromatophore\"",
	"                     /organelle=\"hydrogenosome\"",
	"                     /organelle=\"mitochondrion\"",
	"                     /organelle=\"nucleomorph\"",
	"                     /organelle=\"plastid\"",
	"                     /organelle=\"mitochondrion:kinetoplast\"",
	"                     /organelle=\"plastid:chloroplast\"",
	"                     /organelle=\"plastid:apicoplast\"",
	"                     /organelle=\"plastid:chromoplast\"",
	"                     /organelle=\"plastid:cyanelle\"",
	"                     /organelle=\"plastid:leucoplast\"",
	"                     /organelle=\"plastid:proplastid\"",
	"                     /organism=\"Homo sapiens\"",
	"                     /partial",
	"                     /PCR_conditions=\"Initial denaturation:94degC,1.5min\"",
	"                     /PCR_primers=\"fwd_name: CO1P1, fwd_seq: ttgattttttggtcayccwgaagt,\n                     rev_name: CO1R4, rev_seq: ccwvytardcctarraartgttg\"",
	"                     /PCR_primers=\" fwd_name: hoge1, fwd_seq: cgkgtgtatcttact,\n                     rev_name: hoge2, rev_seq: cg<i>gtgtatcttact\"",
	"                     /PCR_primers=\"fwd_name: CO1P1, fwd_seq: ttgattttttggtcayccwgaagt,\n                     fwd_name: CO1P2, fwd_seq: gatacacaggtcayccwgaagt, rev_name: CO1R4,\n                     rev_seq: ccwvytardcctarraartgttg\"",
	"                     /phenotype=\"erythromycin resistance\"",
	"                     /plasmid=\"C-589\"",
	"                     /pop_variant=\"pop1\"",
	"                     /pop_variant=\"Bear Paw\"",
	"                     /product=\"trypsinogen\"",
	"                     /product=\"trypsin\"",
	"                     /product=\"XYZ neural-specific transcript\"",
	"                     /protein_id=\"AAA12345.1\"",
	"                     /protein_id=\"AAA1234567.1\"",
	"                     /proviral",
	"                     /pseudo",
	"                     /pseudogene=\"processed\"",
	"                     /pseudogene=\"unprocessed\"",
	"                     /pseudogene=\"unitary\"",
	"                     /pseudogene=\"allelic\"",
	"                     /pseudogene=\"unknown\"",
	"                     /rearranged",
	"                     /recombination_class=\"meiotic\"",
	"                     /recombination_class=\"chromosome_breakpoint\"",
	"                     /regulatory_class=\"promoter\"",
	"                     /regulatory_class=\"enhancer\"",
	"                     /regulatory_class=\"ribosome_binding_site\"",
	"                     /replace=\"a\"",
	"                     /replace=\"\"",
	"                     /ribosomal_slippage",
	"                     /rpt_family=\"Alu\"",
	"                     /rpt_type=INVERTED",
	"                     /rpt_unit_range=202..245",
	"                     /rpt_unit_seq=\"aagggc\"",
	"                     /rpt_unit_seq=\"ag(5)tg(8)\"",
	"                     /rpt_unit_seq=\"(AAAGA)6(AAAA)1(AAAGA)12\"",
	"                     /satellite=\"satellite: S1a\"",
	"                     /satellite=\"satellite: alpha\"",
	"                     /satellite=\"satellite: gamma III\"",
	"                     /satellite=\"microsatellite: DC130\"",
	"                     /segment=\"6\"",
	"                     /serotype=\"B1\"",
	"                     /serovar=\"O157:H7\"",
	"                     /sex=\"female\"",
	"                     /sex=\"male\"",
	"                     /sex=\"hermaphrodite\"",
	"                     /sex=\"unisexual\"",
	"                     /sex=\"bisexual\"",
	"                     /sex=\"asexual\"",
	"                     /sex=\"monoecious\"",
	"                     /sex=\"dioecious\"",
	"                     /specimen_voucher=\"UAM:Mamm:52179\"",
	"                     /specimen_voucher=\"AMCC:101706\"",
	"                     /specimen_voucher=\"USNM:field series 8798\"",
	"                     /specimen_voucher=\"personal:Dan Janzen:99-SRNP-2003\"",
	"                     /specimen_voucher=\"99-SRNP-2003\"",
	"                     /standard_name=\"dotted\"",
	"                     /strain=\"BALB/c\"",
	"                     /sub_clone=\"lambda-hIL7.20g\"",
	"                     /submitter_seqid=\"NODE_1\"",
	"                     /sub_species=\"lactis\"",
	"                     /sub_strain=\"abis\"",
	"                     /tag_peptide=90..122",
	"                     /tissue_lib=\"tissue library 772\"",
	"                     /tissue_type=\"liver\"",
	"                     /transgenic",
	"                     /translation=\"MASTFPPWYRGCASTPSLKGLIMCTW\"",
	"                     /transl_except=(pos:213..215,aa:Trp)",
	"                     /transl_except=(pos:1017,aa:TERM)",
	"                     /transl_except=(pos:2000..2001,aa:TERM)",
	"                     /transl_except=(pos:X22222:15..17,aa:Ala)",
	"                     /transl_table=4",
	"                     /trans_splicing",
	"                     /type_material=\"type strain of Escherichia coli\"",
	"                     /type_material=\"holotype of Cercopitheus lomamiensis\"",
	"                     /type_material=\"paratype of Cercopitheus lomamiensis\"",
	"                     /variety=\"insularis\"",
	"                     /calculated_mol_wt=3430",
	"                     /site_type=\"other\"",
	"                     /coded_by=\"NM_000207.3:60..392\"",
	"                     /mutated",
}

func TestQualifierIO(t *testing.T) {
	prefix := strings.Repeat(" ", 21)

	for _, in := range qualifierTests {
		state := pars.FromString(in)
		parser := pars.Exact(QualifierParser(prefix))
		result, err := parser.Parse(state)
		if err != nil {
			t.Errorf("while parsing`\n%s\n`: %v", in, err)
			return
		}
		switch q := result.Value.(type) {
		case QualifierIO:
			b := strings.Builder{}
			n, err := q.Format(prefix).WriteTo(&b)
			if err != nil {
				t.Errorf("qf.WriteTo(w) = %d, %v, want %d, nil", n, err, n)
			}
			out := b.String()
			if out != in {
				t.Errorf("q.Format(%q) = %q, want %q", prefix, out, in)
			}
		default:
			t.Errorf("result.Value.(type) = %T, want %T", q, QualifierIO{})
		}
	}

	if _, err := QualifierParser("").Parse(pars.FromString("/rpt_type=other\n")); err != nil {
		t.Errorf("while parsing `\n/rpt_type=other\n\n`: %v", err)
	}

	for _, in := range []string{
		"/sex?",
		"/sex=female",
		"/number?",
		"/number=4\n/",
		"/pseudo=\"true\"",
	} {
		state := pars.FromString(in)
		parser := pars.Exact(QualifierParser(""))
		_, err := parser.Parse(state)
		if err == nil {
			t.Errorf("while parsing`\n%s\n`: expected error", in)
		}
	}

	in := QualifierIO{"foo", "bar"}
	out := "/foo=\"bar\""

	if in.String() != out {
		t.Errorf("qualifier.String() = %s, want %s", in.String(), out)
	}
}

var featureIOTests = []string{
	`     source          1..465
                     /organism="Homo sapiens"
                     /mol_type="mRNA"
                     /db_xref="taxon:9606"
                     /chromosome="11"
                     /map="11p15.5"
     gene            1..465
                     /gene="INS"
                     /gene_synonym="IDDM; IDDM1; IDDM2; ILPR; IRDN; MODY10"
                     /note="insulin"
                     /db_xref="GeneID:3630"
                     /db_xref="HGNC:HGNC:6081"
                     /db_xref="MIM:176730"
     exon            1..42
                     /gene="INS"
                     /gene_synonym="IDDM; IDDM1; IDDM2; ILPR; IRDN; MODY10"
                     /inference="alignment:Splign:2.1.0"
     exon            43..246
                     /gene="INS"
                     /gene_synonym="IDDM; IDDM1; IDDM2; ILPR; IRDN; MODY10"
                     /inference="alignment:Splign:2.1.0"
     CDS             60..392
                     /gene="INS"
                     /gene_synonym="IDDM; IDDM1; IDDM2; ILPR; IRDN; MODY10"
                     /note="proinsulin; preproinsulin"
                     /codon_start=1
                     /product="insulin preproprotein"
                     /protein_id="NP_000198.1"
                     /db_xref="CCDS:CCDS7729.1"
                     /db_xref="GeneID:3630"
                     /db_xref="HGNC:HGNC:6081"
                     /db_xref="MIM:176730"
                     /translation="MALWMRLLPLLALLALWGPDPAAAFVNQHLCGSHLVEALYLVCG
                     ERGFFYTPKTRREAEDLQVGQVELGGGPGAGSLQPLALEGSLQKRGIVEQCCTSICSL
                     YQLENYCN"
     sig_peptide     60..131
                     /gene="INS"
                     /gene_synonym="IDDM; IDDM1; IDDM2; ILPR; IRDN; MODY10"
                     /inference="COORDINATES: ab initio prediction:SignalP:4.0"
     proprotein      132..389
                     /gene="INS"
                     /gene_synonym="IDDM; IDDM1; IDDM2; ILPR; IRDN; MODY10"
                     /product="proinsulin"
     mat_peptide     132..221
                     /gene="INS"
                     /gene_synonym="IDDM; IDDM1; IDDM2; ILPR; IRDN; MODY10"
                     /product="insulin B chain"
     mat_peptide     228..320
                     /gene="INS"
                     /gene_synonym="IDDM; IDDM1; IDDM2; ILPR; IRDN; MODY10"
                     /product="C-peptide"
     mat_peptide     327..389
                     /gene="INS"
                     /gene_synonym="IDDM; IDDM1; IDDM2; ILPR; IRDN; MODY10"
                     /product="insulin A chain"
     exon            247..465
                     /gene="INS"
                     /gene_synonym="IDDM; IDDM1; IDDM2; ILPR; IRDN; MODY10"
                     /inference="alignment:Splign:2.1.0"`,
	`     tRNA            complement(3838377..3838450)
                     /gene="TRI-GAT1-1"
                     /product="tRNA-Ile"
                     /inference="COORDINATES: profile:tRNAscan-SE:1.23"
                     /note="tRNA-Ile (anticodon GAT) 1-1; Derived by automated
                     computational analysis using gene prediction method:
                     tRNAscan-SE."
                     /anticodon=(pos:complement(3838414..3838416),aa:Ile,
                     seq:gat)
                     /db_xref="GeneID:100189132"
                     /db_xref="HGNC:HGNC:34694"`,
}

func TestFeatureKeylineParser(t *testing.T) {
	parser := pars.Exact(featureKeylineParser("     ", 21))
	for _, in := range []string{
		"     source          ",
		"    source          ",
		"     ",
		"     source",
		"     source 1..39",
		"     source          1..39 ",
	} {
		state := pars.FromString(in)
		if _, err := parser.Parse(state); err == nil {
			t.Errorf("while parsing`\n%s\n`: expected error", in)
		}
	}
}

func TestFeatureIO(t *testing.T) {
	parser := pars.Exact(INSDCTableParser(""))
	for _, in := range featureIOTests {
		state := pars.FromString(in)
		result, err := parser.Parse(state)
		if err != nil {
			t.Errorf("while parsing`\n%s\n`: %v", in, err)
			return
		}
		switch ff := result.Value.(type) {
		case []gts.Feature:
			b := strings.Builder{}
			fmtr := INSDCFormatter{ff, "     ", 21}
			n, err := fmtr.WriteTo(&b)
			if err != nil {
				t.Errorf("f.WriteTo(w) = %d, %v, want %d, nil", n, err, n)
			}
			out := b.String()
			if out != in {
				t.Errorf("f.Format(%q, 21) = %q, want %q", "     ", out, in)
			}
		default:
			t.Errorf("result.Value.(type) = %T, want %T", ff, []gts.Feature{})
		}
	}

	if err := parser(pars.FromString(""), pars.Void); err == nil {
		t.Error("while parsing empty string: expected error")
	}
}
