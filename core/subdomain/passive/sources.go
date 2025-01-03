package passive

import (
	"prismx_cli/core/subdomain/subscraping"
	"prismx_cli/core/subdomain/subscraping/sources/alienvault"
	"prismx_cli/core/subdomain/subscraping/sources/anubis"
	"prismx_cli/core/subdomain/subscraping/sources/archiveis"
	"prismx_cli/core/subdomain/subscraping/sources/commoncrawl"
	"prismx_cli/core/subdomain/subscraping/sources/crtsh"
	"prismx_cli/core/subdomain/subscraping/sources/dnsdumpster"
	"prismx_cli/core/subdomain/subscraping/sources/fofa"
	"prismx_cli/core/subdomain/subscraping/sources/fullhunt"
	"prismx_cli/core/subdomain/subscraping/sources/hackertarget"
	"prismx_cli/core/subdomain/subscraping/sources/hunter"
	"prismx_cli/core/subdomain/subscraping/sources/rapiddns"
	"prismx_cli/core/subdomain/subscraping/sources/riddler"
	"prismx_cli/core/subdomain/subscraping/sources/shodan"
	"prismx_cli/core/subdomain/subscraping/sources/sitedossier"
	"prismx_cli/core/subdomain/subscraping/sources/sonarsearch"
	"prismx_cli/core/subdomain/subscraping/sources/sublist3r"
	"prismx_cli/core/subdomain/subscraping/sources/threatbook"
	"prismx_cli/core/subdomain/subscraping/sources/threatcrowd"
	"prismx_cli/core/subdomain/subscraping/sources/threatminer"
	"prismx_cli/core/subdomain/subscraping/sources/virustotal"
	"prismx_cli/core/subdomain/subscraping/sources/waybackarchive"
	"prismx_cli/core/subdomain/subscraping/sources/zoomeye"
)

// DefaultAllSources contains list of all sources
var DefaultAllSources = []string{
	"alienvault",
	"anubis",
	"archiveis",
	"commoncrawl",
	"crtsh",
	"dnsdumpster",
	"hackertarget",
	"passivetotal",
	"rapiddns",
	"riddler",
	"shodan",
	"sitedossier",
	"sonarsearch",
	"sublist3r",
	"threatbook",
	"threatcrowd",
	"hunter",
	"threatminer",
	"virustotal",
	"waybackarchive",
	"zoomeye",
	"fofa",
	"fullhunt",
}

// Agent is a struct for running passive subdomain enumeration
// against a given host. It wraps subscraping package and provides
// a layer to build upon.
type Agent struct {
	sources map[string]subscraping.Source
}

// New creates a new agent for passive subdomain discovery
func New(sources, exclusions []string) *Agent {
	// Create the agent, insert the sources and remove the excluded sources
	agent := &Agent{sources: make(map[string]subscraping.Source)}

	agent.addSources(sources)
	agent.removeSources(exclusions)

	return agent
}

// addSources adds the given list of sources to the source array
func (a *Agent) addSources(sources []string) {
	for _, source := range sources {
		switch source {
		case "alienvault":
			a.sources[source] = &alienvault.Source{}
		case "anubis":
			a.sources[source] = &anubis.Source{}
		case "archiveis":
			a.sources[source] = &archiveis.Source{}
		case "hunter":
			a.sources[source] = &hunter.Source{}
		case "commoncrawl":
			a.sources[source] = &commoncrawl.Source{}
		case "crtsh":
			a.sources[source] = &crtsh.Source{}
		case "dnsdumpster":
			a.sources[source] = &dnsdumpster.Source{}
		case "hackertarget":
			a.sources[source] = &hackertarget.Source{}
		case "rapiddns":
			a.sources[source] = &rapiddns.Source{}
		case "riddler":
			a.sources[source] = &riddler.Source{}
		case "shodan":
			a.sources[source] = &shodan.Source{}
		case "sitedossier":
			a.sources[source] = &sitedossier.Source{}
		case "sonarsearch":
			a.sources[source] = &sonarsearch.Source{}
		case "sublist3r":
			a.sources[source] = &sublist3r.Source{}
		case "threatbook":
			a.sources[source] = &threatbook.Source{}
		case "threatcrowd":
			a.sources[source] = &threatcrowd.Source{}
		case "threatminer":
			a.sources[source] = &threatminer.Source{}
		case "virustotal":
			a.sources[source] = &virustotal.Source{}
		case "waybackarchive":
			a.sources[source] = &waybackarchive.Source{}
		case "zoomeye":
			a.sources[source] = &zoomeye.Source{}
		case "fofa":
			a.sources[source] = &fofa.Source{}
		case "fullhunt":
			a.sources[source] = &fullhunt.Source{}
		}
	}
}

// removeSources deletes the given sources from the source map
func (a *Agent) removeSources(sources []string) {
	for _, source := range sources {
		delete(a.sources, source)
	}
}
