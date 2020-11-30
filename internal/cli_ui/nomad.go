package cli_ui

import (
	"fmt"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec"
	"github.com/ryanuber/columnize"
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/deployer"
	"sort"
	"strings"
	"time"
)

func NomadFormatDryRun(resp *api.JobPlanResponse, submittedJobSpec *deployer.Spec) (string, error) {
	submittedJob, err := jobspec.Parse(strings.NewReader(submittedJobSpec.Content))
	if err != nil {
		return "", err
	}
	return nomadFormatDryRun(resp, submittedJob), nil
}

// formatDryRun produces a string explaining the results of the dry run.
func nomadFormatDryRun(resp *api.JobPlanResponse, job *api.Job) string {
	var rolling *api.Evaluation
	for _, eval := range resp.CreatedEvals {
		if eval.TriggeredBy == "rolling-update" {
			rolling = eval
		}
	}

	var out string
	if len(resp.FailedTGAllocs) == 0 {
		out = "[bold][green]- All tasks successfully allocated.[reset]\n"
	} else {
		out = formatSystemAllocOutput(resp, job, rolling)
	}

	if rolling != nil {
		out += fmt.Sprintf("[green]- Rolling update, next evaluation will be in %s.\n", rolling.Wait)
	}

	if next := resp.NextPeriodicLaunch; !next.IsZero() && !job.IsParameterized() {
		loc, err := job.Periodic.GetLocation()
		if err != nil {
			out += fmt.Sprintf("[yellow]- Invalid time zone: %v", err)
		} else {
			now := time.Now().In(loc)
			out += fmt.Sprintf("[green]- If submitted now, next periodic launch would be at %s (%s from now).\n",
				formatTime(next), formatTimeDifference(now, next, time.Second))
		}
	}

	out = strings.TrimSuffix(out, "\n")
	return out
}

func formatSystemAllocOutput(resp *api.JobPlanResponse, job *api.Job, rolling *api.Evaluation) (out string) {
	// Change the output depending on if we are a system job or not
	if job.Type != nil && *job.Type == "system" {
		out = "[bold][yellow]- WARNING: Failed to place allocations on all nodes.[reset]\n"
	} else {
		out = "[bold][yellow]- WARNING: Failed to place all allocations.[reset]\n"
	}
	sorted := sortedTaskGroupFromMetrics(resp.FailedTGAllocs)
	for _, tg := range sorted {
		metrics := resp.FailedTGAllocs[tg]

		noun := "allocation"
		if metrics.CoalescedFailures > 0 {
			noun += "s"
		}
		out += fmt.Sprintf("%s[yellow]Task Group %q (failed to place %d %s):\n[reset]", strings.Repeat(" ", 2), tg, metrics.CoalescedFailures+1, noun)
		out += fmt.Sprintf("[yellow]%s[reset]\n\n", formatAllocMetrics(metrics, false, strings.Repeat(" ", 4)))
	}
	if rolling == nil {
		out = strings.TrimSuffix(out, "\n")
	}
	return
}

func formatAllocMetrics(metrics *api.AllocationMetric, scores bool, prefix string) string {
	// Print a helpful message if we have an eligibility problem
	var out string
	if metrics.NodesEvaluated == 0 {
		out += fmt.Sprintf("%s* No nodes were eligible for evaluation\n", prefix)
	}

	// Print a helpful message if the user has asked for a DC that has no
	// available nodes.
	for dc, available := range metrics.NodesAvailable {
		if available == 0 {
			out += fmt.Sprintf("%s* No nodes are available in datacenter %q\n", prefix, dc)
		}
	}

	// Print filter info
	for class, num := range metrics.ClassFiltered {
		out += fmt.Sprintf("%s* Class %q: %d nodes excluded by filter\n", prefix, class, num)
	}
	for cs, num := range metrics.ConstraintFiltered {
		out += fmt.Sprintf("%s* Constraint %q: %d nodes excluded by filter\n", prefix, cs, num)
	}

	// Print exhaustion info
	if ne := metrics.NodesExhausted; ne > 0 {
		out += fmt.Sprintf("%s* Resources exhausted on %d nodes\n", prefix, ne)
	}
	for class, num := range metrics.ClassExhausted {
		out += fmt.Sprintf("%s* Class %q exhausted on %d nodes\n", prefix, class, num)
	}
	for dim, num := range metrics.DimensionExhausted {
		out += fmt.Sprintf("%s* Dimension %q exhausted on %d nodes\n", prefix, dim, num)
	}

	// Print quota info
	for _, dim := range metrics.QuotaExhausted {
		out += fmt.Sprintf("%s* Quota limit hit %q\n", prefix, dim)
	}

	// Print scores
	if scores {
		if len(metrics.ScoreMetaData) > 0 {
			scoreOutput := make([]string, len(metrics.ScoreMetaData)+1)
			var scorerNames []string
			for i, scoreMeta := range metrics.ScoreMetaData {
				// Add header as first row
				if i == 0 {
					scoreOutput[0] = "Node|"

					// sort scores alphabetically
					scores := make([]string, 0, len(scoreMeta.Scores))
					for score := range scoreMeta.Scores {
						scores = append(scores, score)
					}
					sort.Strings(scores)

					// build score header output
					for _, scorerName := range scores {
						scoreOutput[0] += fmt.Sprintf("%v|", scorerName)
						scorerNames = append(scorerNames, scorerName)
					}
					scoreOutput[0] += "final score"
				}
				scoreOutput[i+1] = fmt.Sprintf("%v|", scoreMeta.NodeID)
				for _, scorerName := range scorerNames {
					scoreVal := scoreMeta.Scores[scorerName]
					scoreOutput[i+1] += fmt.Sprintf("%.3g|", scoreVal)
				}
				scoreOutput[i+1] += fmt.Sprintf("%.3g", scoreMeta.NormScore)
			}
			out += formatList(scoreOutput)
		} else {
			// Backwards compatibility for old allocs
			for name, score := range metrics.Scores {
				out += fmt.Sprintf("%s* Score %q = %f\n", prefix, name, score)
			}
		}
	}

	out = strings.TrimSuffix(out, "\n")
	return out
}

func sortedTaskGroupFromMetrics(groups map[string]*api.AllocationMetric) []string {
	tgs := make([]string, 0, len(groups))
	for tg := range groups {
		tgs = append(tgs, tg)
	}
	sort.Strings(tgs)
	return tgs
}

// formatList takes a set of strings and formats them into properly
// aligned output, replacing any blank fields with a placeholder
// for awk-ability.
func formatList(in []string) string {
	columnConf := columnize.DefaultConfig()
	columnConf.Empty = "<none>"
	return columnize.Format(in, columnConf)
}

// formatTime formats the time to string based on RFC822
func formatTime(t time.Time) string {
	if t.Unix() < 1 {
		// It's more confusing to display the UNIX epoch or a zero value than nothing
		return ""
	}
	// Return ISO_8601 time format GH-3806
	return t.Format("2006-01-02T15:04:05Z07:00")
}

// formatTimeDifference takes two times and determines their duration difference
// truncating to a passed unit.
// E.g. formatTimeDifference(first=1m22s33ms, second=1m28s55ms, time.Second) -> 6s
func formatTimeDifference(first, second time.Time, d time.Duration) string {
	return second.Truncate(d).Sub(first.Truncate(d)).String()
}
